package main

import "github.com/danilevy1212/self-updater/internal/models"

func getCurrentServerFileName(am models.ApplicationMeta) string {
	res := "current"

	if am.OS == "windows" {
		res = res + ".exe"
	}

	return res
}

func getNewServerFileName(am models.ApplicationMeta) string {
	res := "new"

	if am.OS == "windows" {
		res = res + ".exe"
	}

	return res
}

// NOTE  Due to time concerns, I'm not doing health checks or rollbacks, however, these
//       could be implemented by simply coping the current binary to a backup file and
//       doing a health check against the new executable (GET /health). If it passes, great,
//       if not, rollback.
// func getBackupServerFileName(am models.ApplicationMeta) string {
// 	res := "backup"
//
// 	if am.OS == "windows" {
// 		res = res + ".exe"
// 	}
//
// 	return res
// }
