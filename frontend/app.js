'use strict';

const LIMIT = 20;
let offset = 0;
let totalCount = 0;

const grid = document.getElementById('pokemon-grid');
const listView = document.getElementById('list-view');
const detailView = document.getElementById('detail-view');
const detailEl = document.getElementById('pokemon-detail');
const prevBtn = document.getElementById('prev-btn');
const nextBtn = document.getElementById('next-btn');
const pageInfo = document.getElementById('page-info');
const searchInput = document.getElementById('search-input');
const searchBtn = document.getElementById('search-btn');
const backBtn = document.getElementById('back-btn');

// ── Lista ────────────────────────────────────────────────────────────────────

async function loadList() {
  grid.innerHTML = '<p class="loading">Cargando...</p>';
  try {
    const data = await window.go.app.App.ListPokemon(offset, LIMIT);
    totalCount = data.Count;
    renderGrid(data.Results);
    updatePagination();
  } catch (err) {
    grid.innerHTML = `<p class="loading" style="color:#e53e3e">${err}</p>`;
  }
}

function renderGrid(items) {
  if (!items || items.length === 0) {
    grid.innerHTML = '<p class="loading">No se encontraron Pokémon.</p>';
    return;
  }
  grid.innerHTML = items.map(item => {
    const id = idFromURL(item.URL);
    const sprite = spriteURL(id);
    return `<div class="card" data-name="${item.Name}">
      <img src="${sprite}" alt="${item.Name}" loading="lazy" />
      <div class="poke-name">${item.Name}</div>
    </div>`;
  }).join('');

  grid.querySelectorAll('.card').forEach(card => {
    card.addEventListener('click', () => showDetail(card.dataset.name));
  });
}

function updatePagination() {
  const page = Math.floor(offset / LIMIT) + 1;
  const pages = Math.ceil(totalCount / LIMIT);
  pageInfo.textContent = `Página ${page} / ${pages}`;
  prevBtn.disabled = offset === 0;
  nextBtn.disabled = offset + LIMIT >= totalCount;
}

prevBtn.addEventListener('click', () => { offset -= LIMIT; loadList(); });
nextBtn.addEventListener('click', () => { offset += LIMIT; loadList(); });

// ── Detalle ──────────────────────────────────────────────────────────────────

async function showDetail(name) {
  listView.classList.add('hidden');
  detailView.classList.remove('hidden');
  detailEl.innerHTML = '<p class="loading">Cargando...</p>';

  try {
    const p = await window.go.app.App.GetPokemon(name);
    renderDetail(p);
  } catch (err) {
    detailEl.innerHTML = `<p class="loading" style="color:#e53e3e">${err}</p>`;
  }
}

function renderDetail(p) {
  const types = (p.Types || []).map(t =>
    `<span class="type-badge type-${t.Name}">${t.Name}</span>`
  ).join('');

  const sprites = `
    <div class="sprites">
      ${p.Sprites.FrontDefault ? `<div><img src="${p.Sprites.FrontDefault}" alt="default"/><span>Normal</span></div>` : ''}
      ${p.Sprites.FrontShiny ? `<div><img src="${p.Sprites.FrontShiny}" alt="shiny"/><span>Shiny</span></div>` : ''}
    </div>`;

  const statsRows = (p.Stats || []).map(s => {
    const pct = Math.round((s.BaseStat / 255) * 100);
    return `<tr>
      <td>${s.Name}</td>
      <td><div class="stat-bar-wrap"><div class="stat-bar" style="width:${pct}%"></div></div></td>
      <td>${s.BaseStat}</td>
    </tr>`;
  }).join('');

  detailEl.innerHTML = `
    <h2>#${p.ID} ${p.Name}</h2>
    ${sprites}
    <div class="types">${types}</div>
    <p class="meta">Altura: ${p.Height / 10} m &nbsp;·&nbsp; Peso: ${p.Weight / 10} kg</p>
    <table class="stats-table">
      <thead><tr><th>Stat</th><th></th><th>Base</th></tr></thead>
      <tbody>${statsRows}</tbody>
    </table>`;
}

backBtn.addEventListener('click', () => {
  detailView.classList.add('hidden');
  listView.classList.remove('hidden');
});

// ── Búsqueda ─────────────────────────────────────────────────────────────────

async function search() {
  const query = searchInput.value.trim().toLowerCase();
  if (!query) { loadList(); return; }
  await showDetail(query);
}

searchBtn.addEventListener('click', search);
searchInput.addEventListener('keydown', e => { if (e.key === 'Enter') search(); });

// ── Helpers ──────────────────────────────────────────────────────────────────

function idFromURL(url) {
  const parts = url.replace(/\/$/, '').split('/');
  return parts[parts.length - 1];
}

function spriteURL(id) {
  return `https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/${id}.png`;
}

// ── Init ─────────────────────────────────────────────────────────────────────

loadList();
