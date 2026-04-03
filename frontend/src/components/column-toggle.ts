// ─── Column visibility toggle ─────────────────────────────────────────────────

import { t } from "../i18n";

export interface ColumnConfig {
  key: string;
  label: string;
  fixed?: boolean;
}

const STORAGE_PREFIX = "column-visibility-";

const ICON_EYE = `<svg class="col-toggle-icon" viewBox="0 0 20 20" fill="currentColor" width="16" height="16"><path d="M10 3C5 3 1.7 7.3.5 10c1.2 2.7 4.5 7 9.5 7s8.3-4.3 9.5-7C18.3 7.3 15 3 10 3zm0 12a5 5 0 110-10 5 5 0 010 10zm0-8a3 3 0 100 6 3 3 0 000-6z"/></svg>`;
const ICON_EYE_OFF = `<svg class="col-toggle-icon" viewBox="0 0 20 20" fill="currentColor" width="16" height="16"><path d="M2.7 1.3L1.3 2.7l3.4 3.4C2.7 8 1 9.8.5 10c1.2 2.7 4.5 7 9.5 7 1.8 0 3.4-.6 4.8-1.5l2.5 2.5 1.4-1.4-16-16zM10 15c-4 0-6.8-3.4-8-5 .6-1 1.7-2.4 3.1-3.5l2.2 2.2a3 3 0 004 4l1.5 1.5c-.9.5-1.8.8-2.8.8zm2.8-4.4L8.4 6.2C8.9 6.1 9.4 6 10 6a4 4 0 014 4c0 .2 0 .4-.1.6zM10 5c4 0 6.8 3.4 8 5-.4.7-1 1.6-1.8 2.5l-1.4-1.4c.5-.6.9-1.2 1.2-1.8-1.3-2-3.6-4.3-6-4.3-.5 0-1 .1-1.5.2L7 3.7c.9-.4 1.9-.7 3-.7z"/></svg>`;

function getHiddenColumns(tableId: string): Set<string> {
  try {
    const raw = localStorage.getItem(STORAGE_PREFIX + tableId);
    if (raw) return new Set(JSON.parse(raw));
  } catch { /* ignore corrupt data */ }
  return new Set();
}

function saveHiddenColumns(tableId: string, hidden: Set<string>): void {
  localStorage.setItem(STORAGE_PREFIX + tableId, JSON.stringify([...hidden]));
}

function applyVisibility(table: HTMLElement, hidden: Set<string>): void {
  table.querySelectorAll<HTMLElement>("[data-col]").forEach((cell) => {
    const col = cell.dataset.col!;
    cell.classList.toggle("col-hidden", hidden.has(col));
  });
}

function updateToggleIcons(table: HTMLElement, hidden: Set<string>): void {
  table.querySelectorAll<HTMLElement>(".col-toggle-btn").forEach((btn) => {
    const col = btn.dataset.toggleCol!;
    const isHidden = hidden.has(col);
    btn.innerHTML = isHidden ? ICON_EYE_OFF : ICON_EYE;
    btn.classList.toggle("col-toggle-btn--hidden", isHidden);
    btn.title = isHidden ? t("columnToggle.show") : t("columnToggle.hide");
  });
}

function renderHiddenTags(
  tagsContainer: HTMLElement,
  hidden: Set<string>,
  columns: ColumnConfig[],
  onRestore: (key: string) => void,
): void {
  tagsContainer.innerHTML = "";
  const hiddenCols = columns.filter((c) => !c.fixed && hidden.has(c.key));

  if (hiddenCols.length === 0) {
    tagsContainer.style.display = "none";
    return;
  }

  tagsContainer.style.display = "flex";
  for (const col of hiddenCols) {
    const tag = document.createElement("button");
    tag.className = "col-hidden-tag";
    tag.type = "button";
    tag.title = t("columnToggle.showNamed", { name: col.label });
    tag.innerHTML = `${col.label} <span class="col-hidden-tag__x">&times;</span>`;
    tag.addEventListener("click", () => onRestore(col.key));
    tagsContainer.appendChild(tag);
  }
}

export function initColumnToggle(tableId: string, columns: ColumnConfig[]): void {
  const table = document.querySelector<HTMLElement>(`[data-table-id="${tableId}"]`);
  if (!table) return;

  const hidden = getHiddenColumns(tableId);

  // Create hidden-tags container above the table (or its scroll wrapper)
  const insertTarget = table.closest(".moves-table-wrap, .abilities-table-wrap, .encounters-table-wrap") || table;
  let tagsContainer = insertTarget.parentElement?.querySelector<HTMLElement>(`.col-hidden-tags[data-for="${tableId}"]`);
  if (!tagsContainer) {
    tagsContainer = document.createElement("div");
    tagsContainer.className = "col-hidden-tags";
    tagsContainer.dataset.for = tableId;
    insertTarget.parentElement?.insertBefore(tagsContainer, insertTarget);
  }

  const refreshTags = () => renderHiddenTags(tagsContainer!, hidden, columns, (key) => {
    hidden.delete(key);
    saveHiddenColumns(tableId, hidden);
    applyVisibility(table, hidden);
    updateToggleIcons(table, hidden);
    refreshTags();
  });

  // Inject toggle buttons into <th> elements
  columns.forEach((col) => {
    if (col.fixed) return;
    const th = table.querySelector<HTMLElement>(`th[data-col="${col.key}"]`);
    if (!th) return;
    // Avoid duplicating buttons on re-init
    if (th.querySelector(".col-toggle-btn")) return;

    const btn = document.createElement("button");
    btn.className = "col-toggle-btn" + (hidden.has(col.key) ? " col-toggle-btn--hidden" : "");
    btn.dataset.toggleCol = col.key;
    btn.type = "button";
    btn.title = hidden.has(col.key) ? t("columnToggle.show") : t("columnToggle.hide");
    btn.innerHTML = hidden.has(col.key) ? ICON_EYE_OFF : ICON_EYE;
    btn.addEventListener("click", (e) => {
      e.stopPropagation();
      if (hidden.has(col.key)) {
        hidden.delete(col.key);
      } else {
        hidden.add(col.key);
      }
      saveHiddenColumns(tableId, hidden);
      applyVisibility(table, hidden);
      updateToggleIcons(table, hidden);
      refreshTags();
    });
    th.appendChild(btn);
  });

  applyVisibility(table, hidden);
  refreshTags();
}

export function reapplyColumnVisibility(tableId: string): void {
  const table = document.querySelector<HTMLElement>(`[data-table-id="${tableId}"]`);
  if (!table) return;
  const hidden = getHiddenColumns(tableId);
  applyVisibility(table, hidden);
}
