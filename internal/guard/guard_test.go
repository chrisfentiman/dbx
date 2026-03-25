package guard

import "testing"

func TestGuard_AllowedQueries(t *testing.T) {
	g := New([]string{"test_schema", "gold_analytics"}, "hive_metastore")

	allowed := []struct {
		name string
		sql  string
	}{
		{"simple select", "SELECT * FROM test_schema.my_table"},
		{"select with where", "SELECT id, name FROM test_schema.users WHERE active = true"},
		{"select count", "SELECT COUNT(*) FROM test_schema.deals"},
		{"select with join", "SELECT a.id, b.name FROM test_schema.t1 a JOIN test_schema.t2 b ON a.id = b.id"},
		{"select with CTE", "WITH cte AS (SELECT id FROM test_schema.t1) SELECT * FROM cte"},
		{"select with subquery", "SELECT * FROM (SELECT id FROM test_schema.t1) sub"},
		{"select with union", "SELECT id FROM test_schema.t1 UNION ALL SELECT id FROM test_schema.t2"},
		{"fully qualified 3-part", "SELECT * FROM hive_metastore.test_schema.my_table"},
		{"gold schema", "SELECT * FROM gold_analytics.report"},
		{"select with group by", "SELECT agent_id, COUNT(*) FROM test_schema.deals GROUP BY agent_id ORDER BY COUNT(*) DESC LIMIT 20"},
		{"bare table name", "SELECT 1"},
	}

	for _, tc := range allowed {
		t.Run(tc.name, func(t *testing.T) {
			result := g.Validate(tc.sql)
			if !result.Allowed {
				t.Errorf("expected ALLOWED for %q, got blocked: %s", tc.sql, result.Reason)
			}
		})
	}
}

func TestGuard_BlockedQueries(t *testing.T) {
	g := New([]string{"test_schema"}, "hive_metastore")

	blocked := []struct {
		name   string
		sql    string
		reason string
	}{
		{"insert", "INSERT INTO test_schema.t1 VALUES (1)", "INSERT"},
		{"update", "UPDATE test_schema.t1 SET x = 1", "UPDATE"},
		{"delete", "DELETE FROM test_schema.t1", "DELETE"},
		{"drop table", "DROP TABLE test_schema.t1", "DROP"},
		{"create table", "CREATE TABLE test_schema.t1 (id INT)", "CREATE"},
		{"alter table", "ALTER TABLE test_schema.t1 ADD COLUMN x INT", "ALTER"},
		{"truncate", "TRUNCATE TABLE test_schema.t1", "TRUNCATE"},
		{"merge", "MERGE INTO test_schema.t1 USING test_schema.t2 ON t1.id = t2.id WHEN MATCHED THEN UPDATE SET x = 1", "MERGE"},
		{"grant", "GRANT SELECT ON test_schema.t1 TO user1", "GRANT"},
		{"wrong schema", "SELECT * FROM unauthorized_schema.secret_table", "not in the allowed list"},
	}

	for _, tc := range blocked {
		t.Run(tc.name, func(t *testing.T) {
			result := g.Validate(tc.sql)
			if result.Allowed {
				t.Errorf("expected BLOCKED for %q, but was allowed", tc.sql)
			}
		})
	}
}

func TestGuard_ParseFailure_DenyByDefault(t *testing.T) {
	g := New([]string{"test_schema"}, "hive_metastore")

	// Databricks-specific syntax that the MySQL parser can't handle
	blocked := []struct {
		name string
		sql  string
	}{
		{"unparseable garbage", "FOOBAR BAZQUX test_schema.t1"},
		{"set variable", "SET spark.sql.shuffle.partitions = 200"},
		{"optimize", "OPTIMIZE test_schema.t1"},
	}

	for _, tc := range blocked {
		t.Run(tc.name, func(t *testing.T) {
			result := g.Validate(tc.sql)
			if result.Allowed {
				t.Errorf("expected BLOCKED for unparseable %q, but was allowed", tc.sql)
			}
		})
	}
}

func TestGuard_ParseFailure_AllowReadOnly(t *testing.T) {
	g := New([]string{"test_schema"}, "hive_metastore")

	// Databricks-specific read-only queries that the parser can't handle
	// but start with known read-only prefixes
	allowed := []struct {
		name string
		sql  string
	}{
		{"describe history", "DESCRIBE HISTORY test_schema.t1"},
		{"show tables", "SHOW TABLES IN test_schema"},
		{"select with lateral", "SELECT * FROM test_schema.t1 LATERAL VIEW EXPLODE(arr) t AS val"},
	}

	for _, tc := range allowed {
		t.Run(tc.name, func(t *testing.T) {
			result := g.Validate(tc.sql)
			if !result.Allowed {
				t.Errorf("expected ALLOWED for read-only %q, got blocked: %s", tc.sql, result.Reason)
			}
		})
	}
}

func TestGuard_NoSchemaRestriction(t *testing.T) {
	g := New(nil, "hive_metastore") // empty allowed schemas = no restriction
	result := g.Validate("SELECT * FROM any_schema.any_table")
	if !result.Allowed {
		t.Errorf("expected ALLOWED when no schema restrictions, got: %s", result.Reason)
	}
}

func TestGuard_KeywordInIdentifier(t *testing.T) {
	g := New([]string{"test_schema"}, "hive_metastore")
	// "updated_at" contains "UPDATE" but should NOT be blocked
	result := g.Validate("SELECT updated_at FROM test_schema.t1")
	if !result.Allowed {
		t.Errorf("expected ALLOWED (keyword in identifier), got: %s", result.Reason)
	}
}
