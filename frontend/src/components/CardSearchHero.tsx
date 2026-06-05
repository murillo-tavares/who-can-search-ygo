import { useEffect, useRef, type ReactNode } from "react";

type CardSearchHeroProps = {
  children: ReactNode;
  isSelected: boolean;
};

export function CardSearchHero({ children, isSelected }: CardSearchHeroProps) {
  const rootRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!rootRef.current || !canUseWebGL()) {
      return;
    }

    let isMounted = true;
    let fog: { destroy: () => void } | undefined;

    async function loadFog() {
      const [{ default: FOG }, THREE] = await Promise.all([import("vanta/dist/vanta.fog.min"), import("three")]);
      if (!isMounted || !rootRef.current) {
        return;
      }

      fog = FOG({
        THREE,
        baseColor: 0x06070a,
        blurFactor: 0.58,
        el: rootRef.current,
        gyroControls: false,
        highlightColor: 0xb3bac4,
        lowlightColor: 0x121419,
        midtoneColor: 0x555c66,
        minHeight: 200,
        minWidth: 200,
        mouseControls: true,
        scale: 1,
        scaleMobile: 1,
        speed: 0.75,
        touchControls: true,
        zoom: 0.74,
      });
    }

    void loadFog();

    return () => {
      isMounted = false;
      fog?.destroy();
    };
  }, []);

  return (
    <div
      className={isSelected ? "card-search-hero card-search-hero-selected" : "card-search-hero"}
      ref={rootRef}
    >
      <div className="cinematic-backdrop" aria-hidden="true">
        <div className="cinematic-noise" />
        <div className="cinematic-vignette" />
      </div>
      <div className="card-search-hero-content">{children}</div>
    </div>
  );
}

function canUseWebGL() {
  if (import.meta.env.MODE === "test") {
    return false;
  }

  try {
    const canvas = document.createElement("canvas");
    return Boolean(canvas.getContext("webgl") || canvas.getContext("experimental-webgl"));
  } catch {
    return false;
  }
}
