// Copyright 2024 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tetratelabs/run"
	"github.com/tetratelabs/telemetry"
	"github.com/tetratelabs/telemetry/scope"
)

func TestLoggingSetup(t *testing.T) {
	l1 := scope.Register("l1", "test logger one")
	l2 := scope.Register("l2", "test logger two")

	tests := []struct {
		levels    string
		l1        telemetry.Level
		l2        telemetry.Level
		expectErr bool
	}{
		{"l1:debug", telemetry.LevelDebug, telemetry.LevelInfo, false},
		{"l1:debug,l2:debug", telemetry.LevelDebug, telemetry.LevelDebug, false},
		{"invalid:debug,l2:error", telemetry.LevelInfo, telemetry.LevelError, false},
		{"all:none,l1:debug", telemetry.LevelNone, telemetry.LevelNone, false},
		{"", telemetry.LevelInfo, telemetry.LevelInfo, true},
		{",", telemetry.LevelInfo, telemetry.LevelInfo, true},
		{":", telemetry.LevelInfo, telemetry.LevelInfo, true},
		{"invalid", telemetry.LevelInfo, telemetry.LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.levels, func(t *testing.T) {
			g := run.Group{Logger: telemetry.NoopLogger()}
			g.Register(NewLogSystem(telemetry.NoopLogger()))
			require.Equal(t, tt.expectErr, g.Run("", "--log-levels", tt.levels) != nil)

			require.Equal(t, tt.l1, l1.Level())
			require.Equal(t, tt.l2, l2.Level())
		})
	}
}
