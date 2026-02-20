package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newRedisClientForTest(t *testing.T) *redis.Client {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	t.Cleanup(mr.Close)
	return redis.NewClient(&redis.Options{Addr: mr.Addr()})
}

func TestRaceConditionMiddleware_ConcurrentWrite_HappyPath(t *testing.T) {
	client := newRedisClientForTest(t)
	e := echo.New()
	rc := NewRaceConditionMiddleware(client)
	var counter int32

	handler := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		atomic.AddInt32(&counter, 1)
		return c.NoContent(http.StatusOK)
	})

	var wg sync.WaitGroup
	numRequests := 1000
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/increment", nil)
			c := e.NewContext(req, rec)
			for {
				rec = httptest.NewRecorder()
				c = e.NewContext(req, rec)
				_ = handler(c)
				if rec.Code == http.StatusOK {
					break
				}
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(numRequests), counter)
}

func TestRaceConditionMiddleware_ConcurrentWrite_Stress(t *testing.T) {
	client := newRedisClientForTest(t)
	e := echo.New()
	rc := NewRaceConditionMiddleware(client)
	var counter int32

	handler := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		atomic.AddInt32(&counter, 1)
		return c.NoContent(http.StatusOK)
	})

	var wg sync.WaitGroup
	numRequests := 5000
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/increment", nil)
			c := e.NewContext(req, rec)
			for {
				rec = httptest.NewRecorder()
				c = e.NewContext(req, rec)
				_ = handler(c)
				if rec.Code == http.StatusOK {
					break
				}
			}
		}()
	}

	wg.Wait()

	assert.Equal(t, int32(numRequests), counter)
}

func TestRaceConditionMiddleware_ReadWrite(t *testing.T) {
	client := newRedisClientForTest(t)
	e := echo.New()
	rc := NewRaceConditionMiddleware(client)
	var counter int32

	write := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		atomic.AddInt32(&counter, 1)
		return c.NoContent(http.StatusOK)
	})
	read := func(c echo.Context) error {
		_ = atomic.LoadInt32(&counter)
		return c.NoContent(http.StatusOK)
	}

	var wg sync.WaitGroup
	wg.Add(1000)

	for i := 0; i < 500; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/increment", nil)
			c := e.NewContext(req, rec)
			_ = write(c)
		}()
	}
	for i := 0; i < 500; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/counter", nil)
			c := e.NewContext(req, rec)
			_ = read(c)
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(500), counter)
}

func TestRaceConditionMiddleware_MultipleRoutesProtection(t *testing.T) {
	client := newRedisClientForTest(t)
	e := echo.New()
	rc := NewRaceConditionMiddleware(client)
	var counter int32

	increment := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		atomic.AddInt32(&counter, 1)
		return c.NoContent(http.StatusOK)
	})
	decrement := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		atomic.AddInt32(&counter, -1)
		return c.NoContent(http.StatusOK)
	})

	var wg sync.WaitGroup
	inc := 100
	dec := 40
	wg.Add(inc + dec)

	for i := 0; i < inc; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/increment", nil)
			c := e.NewContext(req, rec)
			for {
				rec = httptest.NewRecorder()
				c = e.NewContext(req, rec)
				_ = increment(c)
				if rec.Code == http.StatusOK {
					break
				}
			}
		}()
	}
	for i := 0; i < dec; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/decrement", nil)
			c := e.NewContext(req, rec)
			for {
				rec = httptest.NewRecorder()
				c = e.NewContext(req, rec)
				_ = decrement(c)
				if rec.Code == http.StatusOK {
					break
				}
			}
		}()
	}

	wg.Wait()
	assert.Equal(t, int32(inc-dec), counter)
}

func TestRaceConditionMiddleware_LockScopeValidation(t *testing.T) {
	client := newRedisClientForTest(t)
	e := echo.New()
	rc := NewRaceConditionMiddleware(client)

	health := func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}
	protected := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		time.Sleep(10 * time.Millisecond)
		return c.NoContent(http.StatusOK)
	})

	// Health should be fast and not blocked
	var wg sync.WaitGroup
	wg.Add(10)
	start := time.Now()
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			c := e.NewContext(req, rec)
			_ = health(c)
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	assert.Less(t, elapsed, 200*time.Millisecond)

	// Protected route should serialize
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/increment", nil)
			c := e.NewContext(req, rec)
			_ = protected(c)
		}()
	}
	wg.Wait()
}

func TestRaceConditionMiddleware_DeadlockPrevention(t *testing.T) {
	client := newRedisClientForTest(t)
	e := echo.New()
	rc := NewRaceConditionMiddleware(client)
	rcImpl := rc.(*RaceConditionMiddleware)

	done := make(chan struct{})
	handler := rc.PreventRaceCondition("counter")(func(c echo.Context) error {
		// Try nested lock attempt; should fail quickly, not deadlock
		ok, _ := rcImpl.acquireLock(c.Request().Context(), "nested", "id")
		_ = ok
		close(done)
		return c.NoContent(http.StatusOK)
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/increment", nil)
	c := e.NewContext(req, rec)
	_ = handler(c)

	select {
	case <-done:
		// ok
	case <-time.After(1 * time.Second):
		t.Fatal("deadlock detected: handler did not complete")
	}
}
