package blockeddomains

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orion-tec/oriondns/db"
	"github.com/orion-tec/oriondns/internal/testutil"
)

func TestBlockedDomainsDB_Insert(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	blockedDomainsDB := New(database)

	ctx := context.Background()
	domain := "malware.com"
	recursive := false

	err := blockedDomainsDB.Insert(ctx, domain, recursive)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM blocked_domains WHERE domain = $1", domain).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestBlockedDomainsDB_Insert_Recursive(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	blockedDomainsDB := New(database)

	ctx := context.Background()
	domain := ".ads.example.com"
	recursive := true

	err := blockedDomainsDB.Insert(ctx, domain, recursive)
	require.NoError(t, err)

	var storedRecursive bool
	err = pool.QueryRow(ctx, "SELECT recursive FROM blocked_domains WHERE domain = $1", domain).Scan(&storedRecursive)
	require.NoError(t, err)
	assert.True(t, storedRecursive)
}

func TestBlockedDomainsDB_GetAll(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	blockedDomainsDB := New(database)

	ctx := context.Background()

	domains := []struct {
		domain    string
		recursive bool
	}{
		{"malware.com", false},
		{".ads.example.com", true},
		{"phishing.net", false},
	}

	for _, d := range domains {
		err := blockedDomainsDB.Insert(ctx, d.domain, d.recursive)
		require.NoError(t, err)
	}

	results, err := blockedDomainsDB.GetAll(ctx)
	require.NoError(t, err)
	require.Len(t, results, 3)

	domainMap := make(map[string]bool)
	for _, result := range results {
		domainMap[result.Domain] = result.Recursive
	}

	assert.False(t, domainMap["malware.com"])
	assert.True(t, domainMap[".ads.example.com"])
	assert.False(t, domainMap["phishing.net"])
}

func TestBlockedDomainsDB_GetAll_Empty(t *testing.T) {
	testutil.SkipIfNoDatabase(t)
	pool := testutil.SetupTestDB(t)
	testutil.TruncateAllTables(t, pool)

	database := db.New(pool)
	blockedDomainsDB := New(database)

	ctx := context.Background()

	results, err := blockedDomainsDB.GetAll(ctx)
	require.NoError(t, err)
	assert.Empty(t, results)
}