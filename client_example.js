// Sample frontend telemetry script
function sendTelemetry(telemetryData) {
    const url = 'http://localhost:8080/telemetry';

    const payload = {
        host: window.location.hostname,
        package: "com.domain.test",
        full_url: window.location.href,
        data: telemetryData
    };

    fetch(url, {
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
    })
    .then(response => response.json())
    .then(data => console.log('Telemetry sent:', data))
    .catch(error => console.error('Error sending telemetry:', error));
}

// Example usage:
sendTelemetry({
    event: 'page_load',
    load_time_ms: 250,
    user_agent: navigator.userAgent
});
