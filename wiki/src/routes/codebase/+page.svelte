<h1>Codebase</h1>

<p>This section explains how the codebase is structured.</p>

<h2>Structure</h2>
<pre>
v1 API - Go
v1 GQL - Go
v1 Websocket - Go
Client - Rust
Email Sender - Rust
Building Scripts - Python
SQL - MariaDB
Cache - Redis
DNS - Cloudflare
Wiki - SvelteKit / Typescript
Website - SvelteKit / Typescript
</pre>

<h2>Build System</h2>
<p>
Repository workflows that are used to keep up with latest updates and tests include:
</p>
<pre>
Dependabot
Custom Rust Workaround to make the reqqest_unstable feature available to the compiler
</pre>

<h2>Building Client</h2>
<p>
  To build the client, we use our Python scripts for certain operating systems.
</p>
<p2>
  We organize them by Linux, Windows, and macOS. 
</p2>
<pre>
  Linux - Flatpak, Ubuntu (Noble Numbat), Ubuntu (Jammy Jellyfish), Debian, Arch Linux, Fedora, Gentoo
  Windows - exe
  macOS - dmg
</pre>