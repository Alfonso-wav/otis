import * as echarts from "echarts/core";
import { RadarChart } from "echarts/charts";
import { TooltipComponent, RadarComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import { ALL_TYPES, effectiveness, TYPE_COLORS } from "../pages/explore/type-chart";
import { t } from "../i18n";

echarts.use([RadarChart, RadarComponent, TooltipComponent, CanvasRenderer]);

let radarInstance: echarts.ECharts | null = null;

export function renderDefenseRadar(
  container: HTMLElement,
  typeName: string,
): void {
  if (radarInstance) {
    radarInstance.dispose();
  }

  radarInstance = echarts.init(container);

  const color = TYPE_COLORS[typeName] ?? "#718096";

  // For the selected type as defender, show how each attacker fares
  const indicators = ALL_TYPES.map((atk) => ({
    name: t(`typeNames.${atk}`),
    max: 2,
  }));

  const values = ALL_TYPES.map((atk) => effectiveness(atk, typeName));

  radarInstance.setOption({
    tooltip: {
      trigger: "item",
    },
    radar: {
      indicator: indicators,
      radius: "70%",
      center: ["50%", "55%"],
      shape: "polygon",
      axisName: {
        color: "#4a5568",
        fontSize: 11,
      },
      splitArea: {
        areaStyle: {
          color: ["#fff", "#f7fafc"],
        },
      },
      splitLine: {
        lineStyle: {
          color: "#e2e8f0",
        },
      },
    },
    series: [
      {
        type: "radar",
        data: [
          {
            value: values,
            name: t(`typeNames.${typeName}`),
            areaStyle: { color: hexToRgba(color, 0.25) },
            lineStyle: { color, width: 2 },
            itemStyle: { color },
          },
        ],
      },
    ],
  });

  window.addEventListener("resize", () => radarInstance?.resize());
}

export function disposeDefenseRadar(): void {
  if (radarInstance) {
    radarInstance.dispose();
    radarInstance = null;
  }
}

function hexToRgba(hex: string, alpha: number): string {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}
