// Statping
// Copyright (C) 2018.  Hunter Long and the project contributors
// Written by Hunter Long <info@socialeck.com> and the project contributors
//
// https://github.com/statping/statping
//
// The licenses for most software and other practical works are designed
// to take away your freedom to share and change the works.  By contrast,
// the GNU General Public License is intended to guarantee your freedom to
// share and change all versions of a program--to make sure it remains free
// software for all its users.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package handlers

import (
	"github.com/statping/statping/types/core"
	"github.com/statping/statping/types/services"
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if !core.App.Setup {
		http.Redirect(w, r, "/setup", http.StatusSeeOther)
		return
	}
	ExecuteResponse(w, r, "base.gohtml", core.App, nil)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"services": len(services.All()),
		"online":   true,
		"setup":    core.App.Setup,
	}
	returnJson(health, w, r)
}
