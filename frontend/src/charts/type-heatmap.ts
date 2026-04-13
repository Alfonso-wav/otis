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
let resizeHandler: (() => void) | null = null;

const TYPE_COUNT = 18;
const NARROW_BREAKPOINT = 480;

function clamp(v: number, lo: number, hi: number): number {
  return Math.max(lo, Math.min(hi, v));
}

function multiplierLabel(mult: number): string {
  if (mult === 0) return t("typeChart.immune");
  if (mult === 0.5) return t("typeChart.resisted");
  if (mult === 1) return t("typeChart.neutral");
  if (mult === 2) return t("typeChart.superEffective");
  return `${mult}`;
}

interface Layout {
  gridTop: number;
  gridLeft: number;
  gridBottom: number;
  gridRight: number;
  labelFontSize: number;
  cellLabelFontSize: number;
  showCellLabel: boolean;
}

function computeLayout(width: number): Layout {
  const narrow = width < NARROW_BREAKPOINT;
  const gridLeft = narrow ? 56 : 80;
  const gridTop = narrow ? 70 : 80;
  const gridRight = narrow ? 16 : 40;
  const gridBottom = narrow ? 24 : 40;
  // Cell size approx: (width - gridLeft - gridRight) / 18
  const cellPx = Math.max(8, (width - gridLeft - gridRight) / TYPE_COUNT);
  const labelFontSize = clamp(Math.round(cellPx * 0.7), 9, 14);
  const cellLabelFontSize = clamp(Math.round(cellPx * 0.55), 8, 12);
  const showCellLabel = cellPx >= 18;
  return {
    gridTop,
    gridLeft,
    gridBottom,
    gridRight,
    labelFontSize,
    cellLabelFontSize,
    showCellLabel,
  };
}

function buildOption(layout: Layout): echarts.EChartsCoreOption {
  const types = [...ALL_TYPES];
  const typeLabels = types.map((tp) => t(`typeNames.${tp}`));

  const data: [number, number, number][] = [];
  for (let y = 0; y < types.length; y++) {
    for (let x = 0; x < types.length; x++) {
      data.push([x, y, effectiveness(types[y], types[x])]);
    }
  }

  return {
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
      top: layout.gridTop,
      bottom: layout.gridBottom,
      left: layout.gridLeft,
      right: layout.gridRight,
      containLabel: false,
    },
    xAxis: {
      type: "category",
      data: typeLabels,
      position: "top",
      axisLabel: {
        rotate: 45,
        fontSize: layout.labelFontSize,
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
        fontSize: layout.labelFontSize,
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
          show: layout.showCellLabel,
          formatter(params: { data: [number, number, number] }) {
            const val = params.data[2];
            if (val === 2) return "2x";
            if (val === 0.5) return "\u00BD";
            if (val === 0) return "0";
            return "";
          },
          fontSize: layout.cellLabelFontSize,
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
  };
}

export function renderTypeHeatmap(
  container: HTMLElement,
  onTypeSelect: (type: string) => void,
): void {
  if (chartInstance) {
    chartInstance.dispose();
    chartInstance = null;
  }
  if (resizeHandler) {
    window.removeEventListener("resize", resizeHandler);
    resizeHandler = null;
  }

  // Defer init to next frame so container dimensions are resolved by CSS
  requestAnimationFrame(() => {
    chartInstance = echarts.init(container);

    const types = [...ALL_TYPES];

    const initLayout = computeLayout(container.clientWidth || 600);
    chartInstance.setOption(buildOption(initLayout));

    chartInstance.on("click", (params) => {
      const d = (params as { data?: unknown }).data;
      if (Array.isArray(d) && d.length >= 3) {
        const yIdx = d[1] as number;
        onTypeSelect(types[yIdx]);
      }
    });

    // Debounced resize handler recomputes layout + resizes
    let resizeTimer: ReturnType<typeof setTimeout> | null = null;
    resizeHandler = () => {
      if (resizeTimer) clearTimeout(resizeTimer);
      resizeTimer = setTimeout(() => {
        if (!chartInstance) return;
        const layout = computeLayout(container.clientWidth || 600);
        chartInstance.setOption(buildOption(layout));
        chartInstance.resize();
      }, 150);
    };
    window.addEventListener("resize", resizeHandler);
  });
}

export function disposeTypeHeatmap(): void {
  if (resizeHandler) {
    window.removeEventListener("resize", resizeHandler);
    resizeHandler = null;
  }
  if (chartInstance) {
    chartInstance.dispose();
    chartInstance = null;
  }
}
