//go:build integration

package bench

import (
	"context"
	"fmt"
	"testing"

	"github.com/dyxj/bigbackend/internal/sqlgen/bigbackend/public/entity"
	"github.com/dyxj/bigbackend/internal/user/invitation"
	"github.com/dyxj/bigbackend/pkg/logx"
	"github.com/dyxj/bigbackend/pkg/sqldb"
	"github.com/dyxj/bigbackend/pkg/testx"
	"github.com/dyxj/bigbackend/test"
	"github.com/dyxj/bigbackend/test/faker"
	"github.com/stretchr/testify/assert"
)

// Benchmark runtime guide:
//
// Quick development(default benchtime=1s and count=1)
//
//	task bench:run name=BenchmarkUpdaterSQLDBUserInvitation_Multiple
//
// Balanced
//
//	task bench:run name=BenchmarkUpdaterSQLDBUserInvitation_Multiple benchtime=1s count=3
//
// Production quality
//
//	task bench:run name=BenchmarkUpdaterSQLDBUserInvitation_Multiple benchtime=3s count=5
//
// Results analysis:
// After running the benchmark, use the "benchstat" tool to analyze results. Since it is in a single file
// use -col parameter to define columns to compare.
//
// Note that result moved to baseline folder so it won't be overwritten by subsequent runs.
//
//	benchstat -col /method test/bench_baseline/BenchmarkUpdaterSQLDBUserInvitation_Multiple_3s_5.txt
//	for quick comparison of count=1 results, alpha can be set to 1
func BenchmarkUpdaterSQLDBUserInvitation_Multiple(b *testing.B) {
	logger, err := logx.InitLogger()
	if err != nil {
		b.Fatalf("failed to initialize logger: %v", err)
	}
	dbConn := testx.GlobalBenchEnv().DBConn()
	creator := invitation.NewCreatorSQLDB(logger)
	updater := invitation.NewUpdaterSQLDB(logger)
	ctx := b.Context()

	batchSizes := []int{1, 5, 10, 25, 50, 100}

	for _, batchSize := range batchSizes {
		b.Run(fmt.Sprintf("method=Single/BatchSize%d", batchSize), func(b *testing.B) {
			// Pre-create records for all iterations
			totalRecords := b.N * batchSize
			records := createBenchmarkRecords(b, ctx, dbConn, creator, totalRecords)

			// Pre-prepare all batches with updated status
			batches := make([][]entity.UserInvitation, b.N)
			for i := 0; i < b.N; i++ {
				batches[i] = prepareBatch(records, i, batchSize)
			}

			b.Cleanup(func() {
				test.TruncateUserInvitation(dbConn)
			})

			// Reset timer to exclude setup time
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				batch := batches[i]

				for j := range batch {
					_, err := updater.UpdateInvitationTx(ctx, dbConn, batch[j])
					assert.NoError(b, err)
				}
			}
		})

		b.Run(fmt.Sprintf("method=Batch/BatchSize%d", batchSize), func(b *testing.B) {
			// Pre-create records for all iterations
			totalRecords := b.N * batchSize
			records := createBenchmarkRecords(b, ctx, dbConn, creator, totalRecords)

			// Pre-prepare all batches with updated status
			batches := make([][]entity.UserInvitation, b.N)
			for i := 0; i < b.N; i++ {
				batches[i] = prepareBatch(records, i, batchSize)
			}

			b.Cleanup(func() {
				test.TruncateUserInvitation(dbConn)
			})

			// Reset timer to exclude setup time
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				batch := batches[i]

				_, err := updater.BatchUpdateInvitationTx(ctx, dbConn, batch)
				assert.NoError(b, err)
			}
		})
	}
}

// createBenchmarkRecords creates a specified number of user invitation records for benchmarking
func createBenchmarkRecords(
	b *testing.B,
	ctx context.Context,
	dbConn sqldb.Executable,
	creator *invitation.CreatorSQLDB,
	count int,
) []entity.UserInvitation {
	records := make([]entity.UserInvitation, count)
	for i := 0; i < count; i++ {
		inv := faker.UserInvitationEntity()
		inv.Email = inv.Email + fmt.Sprintf("+%d", i) // Ensure unique email
		inv.Status = string(invitation.StatusPending)
		inserted, err := creator.InsertUserInvitation(ctx, dbConn, inv)
		assert.NoError(b, err)
		records[i] = inserted
	}
	return records
}

// prepareBatch extracts a batch of records and updates their status to Accepted
func prepareBatch(records []entity.UserInvitation, iteration int, batchSize int) []entity.UserInvitation {
	start := iteration * batchSize
	end := start + batchSize
	batch := records[start:end]

	// Update status for each record in the batch
	for j := range batch {
		batch[j].Status = string(invitation.StatusAccepted)
	}

	return batch
}
