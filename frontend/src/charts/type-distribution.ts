import * as echarts from "echarts/core";
import { PieChart } from "echarts/charts";
import { TooltipComponent } from "echarts/components";
import { CanvasRenderer } from "echarts/renderers";

echarts.use([PieChart, TooltipComponent, CanvasRenderer]);

const TYPE_COLORS: Record<string, string> = {
  normal: "#a0aec0", fire: "#f6ad55", water: "#63b3ed", grass: "#68d391",
  electric: "#f6e05e", psychic: "#f687b3", ice: "#76e4f7", fighting: "#c05621",
  poison: "#9f7aea", ground: "#d69e2e", flying: "#90cdf4", bug: "#a8e063",
  rock: "#b7791f", ghost: "#553c9a", dragon: "#7f9cf5", dark: "#4a5568",
  steel: "#718096", fairy: "#fbb6ce",
};

export function renderTypeDistributionChart(
  containerId: string,
  data: Record<string, number>,
  onTypeClick?: (typeName: string) => void,
): void {
  const el = document.getElementById(containerId);
  if (!el) return;

  if (!data || Object.keys(data).length === 0) {
    el.innerHTML =
      '<p style="color:#718096;font-size:0.85rem;text-align:center;padding:1rem">Sin datos de distribución</p>';
    return;
  }

  // Defer init so the container has computed dimensions after DOM insertion
  requestAnimationFrame(() => {
    const chart = echarts.init(el);

    const seriesData = Object.entries(data)
      .sort((a, b) => b[1] - a[1])
      .map(([type, count]) => ({
        name: type,
        value: count,
        itemStyle: { color: TYPE_COLORS[type] ?? "#718096" },
      }));

    chart.setOption({
      tooltip: {
        trigger: "item",
        formatter: "{b}: {c} Pokémon ({d}%)",
      },
      series: [
        {
          type: "pie",
          radius: ["35%", "65%"],
          center: ["50%", "50%"],
          data: seriesData,
          label: {
            show: true,
            formatter: "{b}",
            fontSize: 11,
          },
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: "rgba(0,0,0,0.3)",
            },
          },
        },
      ],
    });

    chart.resize();

    if (onTypeClick) {
      chart.on("click", (params: { name?: string }) => {
        if (params.name) onTypeClick(params.name);
      });
      el.style.cursor = "pointer";
    }

    window.addEventListener("resize", () => chart.resize());
  });
}
