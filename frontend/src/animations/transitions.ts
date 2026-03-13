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
  const cards = container.querySelectorAll(".card");
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
