import * as echarts from "echarts/core";
import { HeatmapChart } from "echarts/charts";
import {
  GridComponent,
  TooltipComponent,
  VisualMapPiecewiseComponent,
} from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import { ALL_TYPES, effectiveness } from "../pages/explore/type-chart";
import { t } from "../i18n";

echarts.use([
  HeatmapChart,
  GridComponent,
  TooltipComponent,
  VisualMapPiecewiseComponent,
  CanvasRenderer,
]);

let chartInstance: echarts.ECharts | null = null;

function multiplierLabel(mult: number): string {
  if (mult === 0) return t("typeChart.immune");
  if (mult === 0.5) return t("typeChart.resisted");
  if (mult === 1) return t("typeChart.neutral");
  if (mult === 2) return t("typeChart.superEffective");
  return `${mult}`;
}

export function renderTypeHeatmap(
  container: HTMLElement,
  onTypeSelect: (type: string) => void,
): void {
  if (chartInstance) {
    chartInstance.dispose();
  }

  chartInstance = echarts.init(container);

  const types = [...ALL_TYPES];
  const typeLabels = types.map((tp) => t(`typeNames.${tp}`));

  // Build data: [xIndex (defender), yIndex (attacker), multiplier]
  const data: [number, number, number][] = [];
  for (let y = 0; y < types.length; y++) {
    for (let x = 0; x < types.length; x++) {
      data.push([x, y, effectiveness(types[y], types[x])]);
    }
  }

  chartInstance.setOption({
    tooltip: {
      position: "top",
      formatter(params: { data: [number, number, number] }) {
        const [x, y, val] = params.data;
        const atk = t(`typeNames.${types[y]}`);
        const def = t(`typeNames.${types[x]}`);
        return [
          `<strong>${t("typeChart.attackerType")}:</strong> ${atk}`,
          `<strong>${t("typeChart.defenderType")}:</strong> ${def}`,
          `<strong>${t("typeChart.multiplier")}:</strong> ${val}x (${multiplierLabel(val)})`,
        ].join("<br>");
      },
    },
    grid: {
      top: 80,
      bottom: 40,
      left: 80,
      right: 40,
      containLabel: false,
    },
    xAxis: {
      type: "category",
      data: typeLabels,
      position: "top",
      axisLabel: {
        rotate: 45,
        fontSize: 10,
        color: "#4a5568",
      },
      splitArea: { show: false },
      axisTick: { show: false },
      axisLine: { show: false },
    },
    yAxis: {
      type: "category",
      data: typeLabels,
      inverse: true,
      axisLabel: {
        fontSize: 10,
        color: "#4a5568",
      },
      splitArea: { show: false },
      axisTick: { show: false },
      axisLine: { show: false },
    },
    visualMap: {
      type: "piecewise",
      pieces: [
        { value: 0, label: t("typeChart.immune"), color: "#4a5568" },
        { value: 0.5, label: t("typeChart.resisted"), color: "#fc8181" },
        { value: 1, label: t("typeChart.neutral"), color: "#e2e8f0" },
        { value: 2, label: t("typeChart.superEffective"), color: "#68d391" },
      ],
      orient: "horizontal",
      left: "center",
      top: 0,
      textStyle: {
        color: "#4a5568",
        fontSize: 11,
      },
    },
    series: [
      {
        type: "heatmap",
        data,
        label: {
          show: true,
          formatter(params: { data: [number, number, number] }) {
            const val = params.data[2];
            if (val === 2) return "2x";
            if (val === 0.5) return "\u00BD";
            if (val === 0) return "0";
            return "";
          },
          fontSize: 10,
          color: "#2d3748",
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 6,
            shadowColor: "rgba(0,0,0,0.3)",
          },
        },
        itemStyle: {
          borderWidth: 1,
          borderColor: "#fff",
        },
      },
    ],
  });

  chartInstance.on("click", (params) => {
    const d = (params as { data?: unknown }).data;
    if (Array.isArray(d) && d.length >= 3) {
      const yIdx = d[1] as number;
      onTypeSelect(types[yIdx]);
    }
  });

  window.addEventListener("resize", () => chartInstance?.resize());
}

export function disposeTypeHeatmap(): void {
  if (chartInstance) {
    chartInstance.dispose();
    chartInstance = null;
  }
}
