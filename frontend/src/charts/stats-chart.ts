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

export function renderStatsChart(container: HTMLElement, primary: ChartSeries, secondary?: ChartSeries): void {
  if (chartInstance) {
    chartInstance.dispose();
  }

  chartInstance = echarts.init(container);

  const names = primary.stats.map((s) => s.Name);
  const primaryValues = primary.stats.map((s) => s.BaseStat);

  const seriesData: object[] = [
    {
      value: primaryValues,
      name: primary.label,
      areaStyle: {
        color: hexToRgba(primary.color, 0.2),
      },
      lineStyle: {
        color: primary.color,
        width: 2,
      },
      itemStyle: {
        color: primary.color,
      },
    },
  ];

  if (secondary) {
    const secondaryValues = secondary.stats.map((s) => s.BaseStat);
    seriesData.push({
      value: secondaryValues,
      name: secondary.label,
      areaStyle: {
        color: hexToRgba(secondary.color, 0.2),
      },
      lineStyle: {
        color: secondary.color,
        width: 2,
      },
      itemStyle: {
        color: secondary.color,
      },
    });
  }

  chartInstance.setOption({
    tooltip: {
      trigger: "item",
    },
    radar: {
      indicator: names.map((name) => ({
        name,
        max: 255,
      })),
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
