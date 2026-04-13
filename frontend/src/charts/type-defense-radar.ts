import * as echarts from "echarts/core";
import { RadarChart } from "echarts/charts";
import { TooltipComponent, RadarComponent, LegendComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import { ALL_TYPES, effectiveness, TYPE_COLORS } from "../pages/explore/type-chart";
import { t } from "../i18n";

echarts.use([RadarChart, RadarComponent, TooltipComponent, LegendComponent, CanvasRenderer]);

let radarInstance: echarts.ECharts | null = null;
let resizeHandler: (() => void) | null = null;

export interface RadarOptions {
  showDefense: boolean;
  showOffense: boolean;
}

export function renderTypeRadar(
  container: HTMLElement,
  typeName: string,
  options: RadarOptions = { showDefense: true, showOffense: true },
): void {
  if (radarInstance) {
    radarInstance.dispose();
    radarInstance = null;
  }
  if (resizeHandler) {
    window.removeEventListener("resize", resizeHandler);
    resizeHandler = null;
  }

  radarInstance = echarts.init(container);

  const defenseColor = TYPE_COLORS[typeName] ?? "#718096";
  const offenseColor = "#2d3748";

  const indicators = ALL_TYPES.map((type) => ({
    name: t(`typeNames.${type}`),
    max: 2,
  }));

  const defenseValues = ALL_TYPES.map((atk) => effectiveness(atk, typeName));
  const offenseValues = ALL_TYPES.map((def) => effectiveness(typeName, def));

  const defenseName = t("typeChart.defenseRadar");
  const offenseName = t("typeChart.offensiveRadar");

  const data: Array<Record<string, unknown>> = [];
  if (options.showDefense) {
    data.push({
      value: defenseValues,
      name: defenseName,
      areaStyle: { color: hexToRgba(defenseColor, 0.25) },
      lineStyle: { color: defenseColor, width: 2, type: "solid" },
      itemStyle: { color: defenseColor },
      symbol: "circle",
      symbolSize: 6,
    });
  }
  if (options.showOffense) {
    data.push({
      value: offenseValues,
      name: offenseName,
      areaStyle: { color: hexToRgba(offenseColor, 0.15) },
      lineStyle: { color: offenseColor, width: 2, type: "dashed" },
      itemStyle: { color: offenseColor },
      symbol: "diamond",
      symbolSize: 7,
    });
  }

  radarInstance.setOption({
    tooltip: {
      trigger: "item",
    },
    legend: {
      show: false,
      data: [defenseName, offenseName],
    },
    radar: {
      indicator: indicators,
      radius: "65%",
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
        data,
      },
    ],
  });

  resizeHandler = () => radarInstance?.resize();
  window.addEventListener("resize", resizeHandler);
}

// Backwards-compatible alias
export const renderDefenseRadar = (container: HTMLElement, typeName: string): void => {
  renderTypeRadar(container, typeName);
};

export function disposeDefenseRadar(): void {
  if (radarInstance) {
    radarInstance.dispose();
    radarInstance = null;
  }
  if (resizeHandler) {
    window.removeEventListener("resize", resizeHandler);
    resizeHandler = null;
  }
}

function hexToRgba(hex: string, alpha: number): string {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}
