<h1>cast.onion API</h1>

<p>
The cast.onion API provides access to stations, streaming, applications, and administrative tools.
All routes are prefixed with <code>/v1/</code>.
</p>

<h2>Websocket - Session & Connection</h2>
<ul>
  <li><code>GET /v1/ws</code> — WebSocket connection, returns a session ID</li>
  <li><code>GET /v1/docs</code> — Returns this documentation page</li>
</ul>

<h2>Stations</h2>
<ul>
  <li><code>GET /v1/stations</code> — List all active stations (requires session)</li>
  <li><code>GET /v1/stations/:id</code> — Get station by ID (requires session)</li>
</ul>

<h2>Applications</h2>
<ul>
  <li><code>POST /v1/apply</code> — Submit a hosting application (requires session)</li>
  <li><code>GET /v1/apply/:id</code> — Check application status (requires session)</li>
</ul>

<h2>Station Owner</h2>
<ul>
  <li><code>GET /v1/owner/dashboard</code> — Get station info (requires session + X-Access-Token)</li>
  <li><code>PATCH /v1/owner/station</code> — Update station details (requires session + X-Access-Token)</li>
</ul>

<h2>Streaming</h2>
<ul>
  <li><code>POST /v1/broadcast/:station_id</code> — Stream audio to a station (requires X-Station-Key)</li>
  <li><code>GET /v1/listen/:station_id</code> — Listen to live audio stream (audio/ogg)</li>
  <li><code>GET /v1/status/:station_id</code> — Check if a station is live</li>
</ul>

<h2>Guest Rooms</h2>
<ul>
  <li><code>POST /v1/room/create</code> — Create an invite room (requires session + X-Access-Token)</li>
  <li><code>POST /v1/room/join/:code</code> — Join a room using an invite code</li>
  <li><code>GET /v1/room/:room_id</code> — Get room info and guest list</li>
  <li><code>POST /v1/room/:room_id/stream</code> — Guest audio stream (requires X-Guest-ID)</li>
  <li><code>POST /v1/room/:room_id/mute/:guest_id</code> — Host mutes a guest (requires X-Access-Token)</li>
  <li><code>POST /v1/room/:room_id/selfmute</code> — Guest mutes themselves (requires X-Guest-ID)</li>
</ul>

<h2>Admin</h2>
<p>
All admin routes require <code>X-Admin-Token</code>.
</p>

<ul>
  <li><code>GET /v1/admin/applications</code> — List pending applications</li>
  <li><code>POST /v1/admin/applications/:id/approve</code> — Approve application</li>
  <li><code>POST /v1/admin/applications/:id/deny</code> — Deny application</li>
  <li><code>POST /v1/admin/stations/:id/suspend</code> — Suspend station</li>
  <li><code>POST /v1/admin/stations/:id/revoke</code> — Revoke station</li>
  <li><code>POST /v1/admin/stations/:id/unsuspend</code> — Unsuspend station</li>
</ul>

<h2>GraphQL</h2>
<ul>
  <li><code>POST /v1/graphql</code> — GraphQL endpoint for stations, directory, and applications</li>
</ul>

<h2>CDN</h2>
<ul>
  <li><code>GET localhost:1050/:path</code> — Static file server (cdn-files directory)</li>
</ul>