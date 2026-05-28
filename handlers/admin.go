package handlers

import (
	"context"
	"fmt"
	"net/http"

	"cloakapp-receiver/db"
	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	DB *db.InfluxDB
	KS *db.KeyStore
}

func NewAdminHandler(database *db.InfluxDB, keystore *db.KeyStore) *AdminHandler {
	return &AdminHandler{DB: database, KS: keystore}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>Telemetry Dashboard</title>
    <style>
        body { font-family: sans-serif; margin: 0; background: #f4f4f9; display: flex; }
        .sidebar { width: 200px; background: #333; color: white; min-height: 100vh; padding: 20px; }
        .sidebar a { color: white; text-decoration: none; display: block; padding: 10px; margin-bottom: 5px; border-radius: 4px; }
        .sidebar a:hover, .sidebar a.active { background: #4CAF50; }
        .main { flex: 1; padding: 20px; }
        table { width: 100%; border-collapse: collapse; margin-top: 20px; background: white; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        th, td { padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }
        th { background-color: #4CAF50; color: white; }
        tr:hover { background-color: #f5f5f5; }
        .header { display: flex; justify-content: space-between; align-items: center; }
        .btn { padding: 10px 20px; background: #4CAF50; color: white; border: none; cursor: pointer; border-radius: 4px; }
        .btn-danger { background: #f44336; }
        .hidden { display: none; }
    </style>
</head>
<body>
    <div class="sidebar">
        <h2>Admin</h2>
        <a href="#" onclick="showSection('telemetry')" id="link-telemetry" class="active">Telemetry</a>
        <a href="#" onclick="showSection('api-keys')" id="link-api-keys">API Keys</a>
        <a href="/swagger/index.html" target="_blank">Swagger API</a>
    </div>
    <div class="main">
        <div id="section-telemetry">
            <div class="header">
                <h1>Telemetry Data</h1>
                <button class="btn" onclick="loadTelemetry()">Refresh Data</button>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Time</th>
                        <th>Package</th>
                        <th>Host</th>
                        <th>Full URL</th>
                        <th>IP</th>
                        <th>User Agent</th>
                        <th>Data</th>
                    </tr>
                </thead>
                <tbody id="telemetry-body">
                    <tr><td colspan="7">Loading...</td></tr>
                </tbody>
            </table>
        </div>

        <div id="section-api-keys" class="hidden">
            <div class="header">
                <h1>API Key Management</h1>
                <button class="btn" onclick="generateKey()">Generate New Key</button>
            </div>
            <table>
                <thead>
                    <tr>
                        <th>Label</th>
                        <th>Key</th>
                        <th>Created</th>
                        <th>Action</th>
                    </tr>
                </thead>
                <tbody id="keys-body">
                    <tr><td colspan="4">Loading...</td></tr>
                </tbody>
            </table>
        </div>
    </div>

    <script>
        function showSection(id) {
            document.getElementById('section-telemetry').classList.add('hidden');
            document.getElementById('section-api-keys').classList.add('hidden');
            document.getElementById('link-telemetry').classList.remove('active');
            document.getElementById('link-api-keys').classList.remove('active');
            
            document.getElementById('section-' + id).classList.remove('hidden');
            document.getElementById('link-' + id).classList.add('active');
            
            if (id === 'telemetry') loadTelemetry();
            if (id === 'api-keys') loadKeys();
        }

        function loadTelemetry() {
            fetch('/admin/api/data')
                .then(response => response.json())
                .then(data => {
                    const tbody = document.getElementById('telemetry-body');
                    tbody.innerHTML = '';
                    if (!data || data.length === 0) {
                        tbody.innerHTML = '<tr><td colspan="7">No data found</td></tr>';
                        return;
                    }
                    data.forEach(row => {
                        const tr = document.createElement('tr');
                        tr.innerHTML = '<td>' + new Date(row.time).toLocaleString() + '</td>' +
                            '<td>' + row.package + '</td>' +
                            '<td>' + row.host + '</td>' +
                            '<td>' + row.full_url + '</td>' +
                            '<td>' + row.ip + '</td>' +
                            '<td title="' + row.user_agent + '">' + (row.user_agent ? row.user_agent.substring(0, 30) : '') + '...</td>' +
                            '<td><pre>' + JSON.stringify(row.fields, null, 2) + '</pre></td>';
                        tbody.appendChild(tr);
                    });
                });
        }

        function loadKeys() {
            fetch('/admin/api/keys')
                .then(response => response.json())
                .then(data => {
                    const tbody = document.getElementById('keys-body');
                    tbody.innerHTML = '';
                    if (!data || data.length === 0) {
                        tbody.innerHTML = '<tr><td colspan="4">No API keys found</td></tr>';
                        return;
                    }
                    data.forEach(k => {
                        const tr = document.createElement('tr');
                        tr.innerHTML = '<td>' + k.label + '</td>' +
                            '<td><code>' + k.key + '</code></td>' +
                            '<td>' + new Date(k.created).toLocaleString() + '</td>' +
                            '<td><button class="btn btn-danger" onclick="deleteKey(\'' + k.key + '\')">Delete</button></td>';
                        tbody.appendChild(tr);
                    });
                });
        }

        function generateKey() {
            const label = prompt("Enter label for new API key:");
            if (!label) return;
            fetch('/admin/api/keys', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({label: label})
            }).then(() => loadKeys());
        }

        function deleteKey(key) {
            if (!confirm("Are you sure you want to delete this API key?")) return;
            fetch('/admin/api/keys/' + key, { method: 'DELETE' })
                .then(() => loadKeys());
        }

        loadTelemetry();
    </script>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

type TelemetryRow struct {
	Time      string                 `json:"time"`
	Host      string                 `json:"host"`
	Package   string                 `json:"package"`
	FullURL   string                 `json:"full_url"`
	IP        string                 `json:"ip"`
	UserAgent string                 `json:"user_agent"`
	Fields    map[string]interface{} `json:"fields"`
}

// @Summary Get telemetry data
// @Description Retrieve the latest telemetry data from InfluxDB
// @Tags Admin
// @Produce json
// @Success 200 {array} TelemetryRow
// @Router /admin/api/data [get]
func (h *AdminHandler) GetData(c *gin.Context) {
	if h.DB == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not initialized"})
		return
	}

	query := fmt.Sprintf(`from(bucket: "%s")
		|> range(start: -24h)
		|> filter(fn: (r) => r["_measurement"] == "telemetry")
		|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
		|> limit(n: 50)
		|> sort(columns: ["_time"], desc: true)`, h.DB.Bucket)

	result, err := h.DB.QueryAPI.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []TelemetryRow
	for result.Next() {
		values := result.Record().Values()
		
		row := TelemetryRow{
			Time:      result.Record().Time().String(),
			Host:      fmt.Sprintf("%v", values["host"]),
			Package:   fmt.Sprintf("%v", values["package"]),
			FullURL:   fmt.Sprintf("%v", values["full_url"]),
			IP:        fmt.Sprintf("%v", values["ip"]),
			UserAgent: fmt.Sprintf("%v", values["user_agent"]),
			Fields:    make(map[string]interface{}),
		}

		for k, v := range values {
			if k != "_time" && k != "_start" && k != "_stop" && k != "_measurement" && 
			   k != "host" && k != "package" && k != "full_url" && k != "ip" && k != "user_agent" {
				row.Fields[k] = v
			}
		}
		data = append(data, row)
	}

	if result.Err() != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Err().Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary List API keys
// @Tags API Keys
// @Produce json
// @Success 200 {array} db.APIKey
// @Router /admin/api/keys [get]
func (h *AdminHandler) ListKeys(c *gin.Context) {
	c.JSON(http.StatusOK, h.KS.List())
}

type CreateKeyRequest struct {
	Label string `json:"label" binding:"required"`
}

// @Summary Generate API key
// @Tags API Keys
// @Accept json
// @Produce json
// @Param request body CreateKeyRequest true "Key Label"
// @Success 200 {object} map[string]string
// @Router /admin/api/keys [post]
func (h *AdminHandler) CreateKey(c *gin.Context) {
	var req CreateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := h.KS.Generate(req.Label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": key})
}

// @Summary Delete API key
// @Tags API Keys
// @Param key path string true "API Key"
// @Success 200 {object} map[string]string
// @Router /admin/api/keys/{key} [delete]
func (h *AdminHandler) DeleteKey(c *gin.Context) {
	key := c.Param("key")
	if err := h.KS.Delete(key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
