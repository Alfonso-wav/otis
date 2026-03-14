const MAX_SUGGESTIONS = 8;

export function createAutocomplete(
  input: HTMLInputElement,
  names: string[],
  onSelect: (name: string) => void,
): void {
  // Ensure the input's parent is positioned so the dropdown aligns correctly
  const wrapper = input.parentElement!;
  if (getComputedStyle(wrapper).position === "static") {
    wrapper.style.position = "relative";
  }

  let dropdown: HTMLUListElement | null = null;
  let activeIndex = -1;

  function getDropdown(): HTMLUListElement {
    if (!dropdown) {
      dropdown = document.createElement("ul");
      dropdown.className = "autocomplete-dropdown";
      wrapper.appendChild(dropdown);
    }
    return dropdown;
  }

  function closeDropdown(): void {
    dropdown?.remove();
    dropdown = null;
    activeIndex = -1;
  }

  function setActive(index: number): void {
    if (!dropdown) return;
    const items = dropdown.querySelectorAll<HTMLLIElement>("li");
    items.forEach((li, i) => li.classList.toggle("autocomplete-dropdown__item--active", i === index));
    activeIndex = index;
  }

  function openDropdown(query: string): void {
    const q = query.toLowerCase();
    const matches = names.filter((n) => n.toLowerCase().includes(q)).slice(0, MAX_SUGGESTIONS);

    if (matches.length === 0) {
      closeDropdown();
      return;
    }

    const ul = getDropdown();
    ul.innerHTML = matches
      .map((name) => `<li class="autocomplete-dropdown__item" data-name="${name}">${name}</li>`)
      .join("");
    activeIndex = -1;

    ul.querySelectorAll<HTMLLIElement>("li").forEach((li) => {
      li.addEventListener("mousedown", (e) => {
        e.preventDefault(); // prevent input blur before click registers
        select(li.dataset.name!);
      });
    });
  }

  function select(name: string): void {
    input.value = name;
    closeDropdown();
    onSelect(name);
  }

  input.addEventListener("input", () => {
    const v = input.value.trim();
    if (v.length === 0) {
      closeDropdown();
    } else {
      openDropdown(v);
    }
  });

  input.addEventListener("keydown", (e) => {
    if (!dropdown) return;
    const items = dropdown.querySelectorAll<HTMLLIElement>("li");
    if (e.key === "ArrowDown") {
      e.preventDefault();
      setActive(Math.min(activeIndex + 1, items.length - 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActive(Math.max(activeIndex - 1, 0));
    } else if (e.key === "Enter") {
      if (activeIndex >= 0 && items[activeIndex]) {
        e.preventDefault();
        e.stopPropagation();
        select(items[activeIndex].dataset.name!);
      }
    } else if (e.key === "Escape") {
      closeDropdown();
    }
  });

  input.addEventListener("blur", () => {
    // Small delay so mousedown on a list item fires before blur closes it
    setTimeout(closeDropdown, 150);
  });
}
