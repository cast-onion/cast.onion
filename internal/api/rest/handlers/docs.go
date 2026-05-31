package handlers

import (
	"fmt"
	"net/http"
)

func Docs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html><html><head><meta charset="utf-8">
<title>cast.onion api — v1</title>
<link rel="icon" type="image/x-icon" href="http://localhost:1050/web/favicon.ico">   
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#0a0a0a;color:#e0e0e0;font-family:'Courier New',monospace;font-size:13px;padding:48px 32px;max-width:860px}
h1{color:#c8ff00;font-size:18px;font-weight:normal;margin-bottom:8px}
.sub{color:#555;font-size:12px;margin-bottom:40px}
h2{color:#888;font-size:11px;letter-spacing:.08em;text-transform:uppercase;margin:32px 0 12px}
.route{display:grid;grid-template-columns:80px 280px 1fr;gap:0;margin-bottom:2px;padding:10px 14px;background:#111;border-radius:3px}
.method{font-weight:bold}
.get{color:#c8ff00}.post{color:#44cc88}.patch{color:#ffaa00}.delete{color:#ff4444}
.path{color:#e0e0e0}
.desc{color:#555}
.auth{color:#333;font-size:11px;margin-top:2px}
.badge{display:inline-block;padding:1px 6px;border-radius:2px;font-size:10px;margin-left:6px;border:1px solid}
.s{color:#555;border-color:#333}.a{color:#666;border-color:#444}
</style></head><body>
<h1>cast.onion api</h1>
<div class="sub">v1 · api.castonion.xyz · all routes prefixed /v1/</div>

<h2>session &amp; connection</h2>
<div class="route"><span class="method get">GET</span><span class="path">/v1/ws</span><span class="desc">WebSocket — receive session ID on connect</span></div>
<div class="route"><span class="method get">GET</span><span class="path">/v1/docs</span><span class="desc">this page</span></div>

<h2>stations</h2>
<div class="route"><span class="method get">GET</span><span class="path">/v1/stations</span><span class="desc">list all active stations<span class="badge s">session</span></span></div>
<div class="route"><span class="method get">GET</span><span class="path">/v1/stations/{id}</span><span class="desc">get station by id<span class="badge s">session</span></span></div>

<h2>applications</h2>
<div class="route"><span class="method post">POST</span><span class="path">/v1/apply</span><span class="desc">submit a hosting application<span class="badge s">session</span></span></div>
<div class="route"><span class="method get">GET</span><span class="path">/v1/apply/{id}</span><span class="desc">get application status<span class="badge s">session</span></span></div>

<h2>station owner</h2>
<div class="route"><span class="method get">GET</span><span class="path">/v1/owner/dashboard</span><span class="desc">get station info<span class="badge s">session</span> X-Access-Token</span></div>
<div class="route"><span class="method patch">PATCH</span><span class="path">/v1/owner/station</span><span class="desc">update station details<span class="badge s">session</span> X-Access-Token</span></div>

<h2>streaming</h2>
<div class="route"><span class="method post">POST</span><span class="path">/v1/broadcast/{station_id}</span><span class="desc">stream audio to station — X-Station-Key header</span></div>
<div class="route"><span class="method get">GET</span><span class="path">/v1/listen/{station_id}</span><span class="desc">receive live audio stream (audio/ogg)</span></div>
<div class="route"><span class="method get">GET</span><span class="path">/v1/status/{station_id}</span><span class="desc">check if station is live</span></div>

<h2>guest rooms</h2>
<div class="route"><span class="method post">POST</span><span class="path">/v1/room/create</span><span class="desc">create invite room<span class="badge s">session</span> X-Access-Token</span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/room/join/{code}</span><span class="desc">join a room by code<span class="badge s">session</span></span></div>
<div class="route"><span class="method get">GET</span><span class="path">/v1/room/{room_id}</span><span class="desc">get room info and guest list</span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/room/{room_id}/stream</span><span class="desc">guest audio stream — X-Guest-ID header</span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/room/{room_id}/mute/{guest_id}</span><span class="desc">host mutes a guest — X-Access-Token</span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/room/{room_id}/selfmute</span><span class="desc">guest mutes self — X-Guest-ID</span></div>

<h2>admin (X-Admin-Token required)</h2>
<div class="route"><span class="method get">GET</span><span class="path">/v1/admin/applications</span><span class="desc">list pending applications<span class="badge a">admin</span></span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/admin/applications/{id}/approve</span><span class="desc">approve application<span class="badge a">admin</span></span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/admin/applications/{id}/deny</span><span class="desc">deny application<span class="badge a">admin</span></span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/admin/stations/{id}/suspend</span><span class="desc">suspend station<span class="badge a">admin</span></span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/admin/stations/{id}/revoke</span><span class="desc">revoke station<span class="badge a">admin</span></span></div>
<div class="route"><span class="method post">POST</span><span class="path">/v1/admin/stations/{id}/unsuspend</span><span class="desc">unsuspend station<span class="badge a">admin</span></span></div>

<h2>graphql</h2>
<div class="route"><span class="method post">POST</span><span class="path">/v1/graphql</span><span class="desc">GraphQL endpoint — directory, station info, application status<span class="badge s">session</span></span></div>

<h2>cdn</h2>
<div class="route"><span class="method get">GET</span><span class="path">api.castonion.xyz/{path}</span><span class="desc">serve files from cdn-files/ directory</span></div>
</body></html>`)
}
