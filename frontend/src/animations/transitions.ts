import gsap from "gsap";

export function showView(viewIn: HTMLElement, viewOut: HTMLElement): Promise<void> {
  return new Promise((resolve) => {
    const tl = gsap.timeline({
      onComplete: resolve,
    });

    tl.to(viewOut, {
      opacity: 0,
      duration: 0.2,
      ease: "power2.in",
      onComplete() {
        viewOut.style.display = "none";
      },
    });

    tl.fromTo(
      viewIn,
      { opacity: 0, y: 20 },
      {
        opacity: 1,
        y: 0,
        duration: 0.3,
        ease: "power2.out",
        onStart() {
          viewIn.style.display = "";
          viewIn.style.removeProperty("display");
          viewIn.classList.remove("hidden");
        },
      },
    );
  });
}

export function staggerCards(container: HTMLElement): void {
  const cards = container.querySelectorAll(".poke-card");
  if (cards.length === 0) return;

  gsap.fromTo(
    cards,
    { opacity: 0, y: 15 },
    {
      opacity: 1,
      y: 0,
      duration: 0.25,
      stagger: 0.03,
      ease: "power2.out",
    },
  );
}

export function morphToTable(
  container: HTMLElement,
  renderFn: () => Promise<void>,
): Promise<void> {
  return new Promise((resolve) => {
    const children = container.children;
    if (children.length === 0) {
      renderFn().then(resolve);
      return;
    }

    gsap.to(children, {
      opacity: 0,
      scale: 0.85,
      y: -10,
      stagger: 0.015,
      duration: 0.18,
      ease: "power2.in",
      onComplete() {
        renderFn().then(resolve);
      },
    });
  });
}

export function morphToGrid(
  container: HTMLElement,
  renderFn: () => void,
): Promise<void> {
  return new Promise((resolve) => {
    const rows = container.querySelectorAll(".poke-table__row");
    if (rows.length === 0) {
      renderFn();
      resolve();
      return;
    }

    gsap.to(rows, {
      opacity: 0,
      x: 20,
      stagger: 0.015,
      duration: 0.15,
      ease: "power2.in",
      onComplete() {
        renderFn();
        resolve();
      },
    });
  });
}
