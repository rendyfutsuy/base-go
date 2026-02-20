package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/utils"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) {
	if utils.ConfigVars == nil {
		utils.InitConfig("config.json")
	}
}

func TestThrottleMiddleware_WithinLimit(t *testing.T) {
	setup(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Set limit to 10
	utils.ConfigVars.Set("throttle.limit", 10)
	throttleMiddleware := NewThrottleMiddleware()
	h := throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	for i := 0; i < 10; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := h(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestThrottleMiddleware_ExceedLimit(t *testing.T) {
	setup(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Set limit to 5 for testing
	utils.ConfigVars.Set("throttle.limit", 5)
	throttleMiddleware := NewThrottleMiddleware()
	h := throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	for i := 0; i < 10; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := h(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// This one should be rejected
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestThrottleMiddleware_ConcurrentBurst(t *testing.T) {
	setup(t)
	e := echo.New()

	limit := 50
	utils.ConfigVars.Set("throttle.limit", limit)
	throttleMiddleware := NewThrottleMiddleware()
	h := throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	var wg sync.WaitGroup
	var successCount int32
	var failCount int32

	numRequests := 100
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h(c)
			assert.NoError(t, err)

			if rec.Code == http.StatusOK {
				atomic.AddInt32(&successCount, 1)
			} else if rec.Code == http.StatusTooManyRequests {
				atomic.AddInt32(&failCount, 1)
			}
		}()
	}

	wg.Wait()

	// With burst = limit*2, total successes should be <= limit*2
	assert.LessOrEqual(t, int(successCount), limit*2)
}

func TestThrottleMiddleware_ResetAfterWindow(t *testing.T) {
	setup(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Set limit to 2 for testing
	utils.ConfigVars.Set("throttle.limit", 2)
	throttleMiddleware := NewThrottleMiddleware()
	h := throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// Consume the limit
	for i := 0; i < 2; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := h(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// This one should be rejected
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)

	// Simulate window reset by recreating middleware (new store instance)
	throttleMiddleware = NewThrottleMiddleware()
	h = throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	// This one should be accepted
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestThrottleMiddleware_ConfigOverrideEnv(t *testing.T) {
	setup(t)
	os.Setenv("THROTTLE__LIMIT", "50")
	utils.InitConfig("config.json")
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	throttleMiddleware := NewThrottleMiddleware()
	h := throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	for i := 0; i < 100; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		err := h(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}

	// 101st should be 429 (burst = limit*2)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := h(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}

func TestThrottleMiddleware_ConfigJsonDefault(t *testing.T) {
	// Clear env override to use config.json default
	os.Unsetenv("THROTTLE__LIMIT")
	utils.InitConfig("config.json")
	// Verify default value from config.json
	assert.Equal(t, 1000, utils.ConfigVars.Int("throttle.limit"))
}

func TestThrottleMiddleware_PerIPThrottle(t *testing.T) {
	setup(t)
	utils.ConfigVars.Set("throttle.limit", 3)
	throttleMiddleware := NewThrottleMiddleware()
	e := echo.New()
	h := throttleMiddleware.Throttle()(func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// IP A
	reqA := httptest.NewRequest(http.MethodGet, "/", nil)
	reqA.Header.Set("X-Real-IP", "1.1.1.1")
	for i := 0; i < 6; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(reqA, rec)
		err := h(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	recA := httptest.NewRecorder()
	cA := e.NewContext(reqA, recA)
	_ = h(cA)
	assert.Equal(t, http.StatusTooManyRequests, recA.Code)

	// IP B should have independent quota
	reqB := httptest.NewRequest(http.MethodGet, "/", nil)
	reqB.Header.Set("X-Real-IP", "2.2.2.2")
	for i := 0; i < 6; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(reqB, rec)
		err := h(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
	recB := httptest.NewRecorder()
	cB := e.NewContext(reqB, recB)
	_ = h(cB)
	assert.Equal(t, http.StatusTooManyRequests, recB.Code)
}

func TestMiddlewareOrder_ThrottleBeforeRace(t *testing.T) {
	setup(t)
	utils.ConfigVars.Set("throttle.limit", 5)
	throttle := NewThrottleMiddleware()

	// Use miniredis-like scenario not required; race condition with nil redis acts as pass-through
	var counter int32
	race := NewRaceConditionMiddleware(nil)

	e := echo.New()
	handler := throttle.Throttle()(race.PreventRaceCondition("order")(func(c echo.Context) error {
		atomic.AddInt32(&counter, 1)
		return c.String(http.StatusOK, "ok")
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	for i := 0; i < 10; i++ {
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		_ = handler(c)
	}
	// 6th should be blocked by throttle and not increment counter
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	_ = handler(c)

	assert.Equal(t, int32(10), counter)
	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
}
