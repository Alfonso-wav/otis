import * as echarts from "echarts/core";
import { RadarChart } from "echarts/charts";
import { TooltipComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import type { Stat } from "../types";

echarts.use([RadarChart, TooltipComponent, CanvasRenderer]);

let chartInstance: echarts.ECharts | null = null;

export interface ChartSeries {
  label: string;
  stats: Stat[];
  color: string;
}

export function renderStatsChart(container: HTMLElement, series: ChartSeries[]): void {
  if (chartInstance) {
    chartInstance.dispose();
  }

  chartInstance = echarts.init(container);

  const names = series[0]?.stats.map((s) => s.Name) ?? [];

  const seriesData: object[] = series.map((s) => ({
    value: s.stats.map((stat) => stat.BaseStat),
    name: s.label,
    areaStyle: { color: hexToRgba(s.color, 0.2) },
    lineStyle: { color: s.color, width: 2 },
    itemStyle: { color: s.color },
  }));

  chartInstance.setOption({
    tooltip: {
      trigger: "item",
    },
    radar: {
      indicator: names.map((name) => ({
        name,
        max: 255,
      })),
      radius: "80%",
      center: ["50%", "50%"],
      shape: "polygon",
      axisName: {
        color: "#4a5568",
        fontSize: 12,
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
        data: seriesData,
      },
    ],
  });
}

export function disposeChart(): void {
  if (chartInstance) {
    chartInstance.dispose();
    chartInstance = null;
  }
}

function hexToRgba(hex: string, alpha: number): string {
  const r = parseInt(hex.slice(1, 3), 16);
  const g = parseInt(hex.slice(3, 5), 16);
  const b = parseInt(hex.slice(5, 7), 16);
  return `rgba(${r}, ${g}, ${b}, ${alpha})`;
}
