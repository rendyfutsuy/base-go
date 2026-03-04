-- Seed Parameters for Post Levels, Languages, and Topics (with parent relationships)
INSERT INTO "parameters" (
    "id",
    "code",
    "name",
    "value",
    "type",
    "description",
    "created_at",
    "updated_at",
    "deleted_at"
)
VALUES
-- lang
    (uuid_generate_v7(), 'LANG_ID', 'ID', 'id', 'lang', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'LANG_EN', 'EN', 'en', 'lang', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),

-- topic (parent-child will be set below)
    (uuid_generate_v7(), 'TOPIC_3D_MODELLING', '3D Modelling', '3d_modelling', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_PROGRAMMING', 'Programming', 'programming', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_LARAVEL', 'Laravel', 'laravel', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_GO_LANG', 'GO-Lang', 'go_lang', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_PHP', 'PHP', 'php', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_NODE_JS', 'Node.JS', 'node_js', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_FULL_STACK', 'Full Stack', 'full_stack', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_BACKEND', 'Backend', 'backend', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL),
    (uuid_generate_v7(), 'TOPIC_FRONTEND', 'Frontend', 'frontend', 'topic', NULL, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL)
ON CONFLICT (code) DO NOTHING;

-- Set parent relationships for topics (children -> Programming)
UPDATE "parameters" AS child
SET parent_id = parent.id,
    updated_at = CURRENT_TIMESTAMP
FROM (
    SELECT id
    FROM "parameters"
    WHERE "type" = 'topic'
      AND "name" = 'Programming'
      AND "deleted_at" IS NULL
    LIMIT 1
) AS parent
WHERE child."type" = 'topic'
  AND child."name" IN ('Laravel', 'GO-Lang', 'PHP', 'Node.JS', 'Full Stack', 'Backend', 'Frontend')
  AND child."deleted_at" IS NULL
  AND (child."parent_id" IS NULL OR child."parent_id" = uuid_nil());
