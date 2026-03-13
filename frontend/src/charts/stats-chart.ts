import * as echarts from "echarts/core";
import { RadarChart } from "echarts/charts";
import { TooltipComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";
import type { Stat } from "../types";

echarts.use([RadarChart, TooltipComponent, CanvasRenderer]);

let chartInstance: echarts.ECharts | null = null;

export function renderStatsChart(container: HTMLElement, stats: Stat[]): void {
  if (chartInstance) {
    chartInstance.dispose();
  }

  chartInstance = echarts.init(container);

  const names = stats.map((s) => s.Name);
  const values = stats.map((s) => s.BaseStat);

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
        data: [
          {
            value: values,
            name: "Base Stats",
            areaStyle: {
              color: "rgba(229, 62, 62, 0.2)",
            },
            lineStyle: {
              color: "#e53e3e",
              width: 2,
            },
            itemStyle: {
              color: "#e53e3e",
            },
          },
        ],
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
