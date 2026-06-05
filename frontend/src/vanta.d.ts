declare module "vanta/dist/vanta.fog.min" {
  import type * as THREE from "three";

  type VantaFogOptions = {
    THREE: typeof THREE;
    baseColor?: number;
    blurFactor?: number;
    el: HTMLElement;
    gyroControls?: boolean;
    highlightColor?: number;
    lowlightColor?: number;
    midtoneColor?: number;
    minHeight?: number;
    minWidth?: number;
    mouseControls?: boolean;
    scale?: number;
    scaleMobile?: number;
    speed?: number;
    touchControls?: boolean;
    zoom?: number;
  };

  type VantaEffect = {
    destroy: () => void;
  };

  export default function FOG(options: VantaFogOptions): VantaEffect;
}
