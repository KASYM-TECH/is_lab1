package main

import "os"

func EnvOrDefault(env, def string) string {
	if env, ok := os.LookupEnv(env); ok {
		return env
	}
	return def
}
