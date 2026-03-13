import * as echarts from "echarts";

// Pokémon por región (Gen representativa para distribución de tipos ilustrativa)
const REGION_TYPE_DATA: Record<string, Record<string, number>> = {
  kanto: {
    normal: 22, fire: 12, water: 32, grass: 14, electric: 9,
    psychic: 14, ice: 5, fighting: 8, poison: 28, ground: 14,
    flying: 19, bug: 12, rock: 11, ghost: 3, dragon: 3,
  },
  johto: {
    normal: 18, fire: 9, water: 18, grass: 10, electric: 7,
    psychic: 12, ice: 6, fighting: 5, poison: 6, ground: 9,
    flying: 14, bug: 10, rock: 7, ghost: 4, dragon: 5, dark: 7, steel: 8,
  },
  hoenn: {
    normal: 19, fire: 10, water: 28, grass: 15, electric: 7,
    psychic: 11, ice: 8, fighting: 7, poison: 7, ground: 12,
    flying: 16, bug: 13, rock: 12, ghost: 6, dragon: 8, dark: 8, steel: 10,
  },
  sinnoh: {
    normal: 15, fire: 8, water: 16, grass: 12, electric: 8,
    psychic: 15, ice: 9, fighting: 10, ground: 10, flying: 13,
    bug: 9, rock: 8, ghost: 9, dragon: 7, dark: 7, steel: 12, poison: 5,
  },
  unova: {
    normal: 26, fire: 12, water: 17, grass: 18, electric: 10,
    psychic: 16, ice: 8, fighting: 12, poison: 8, ground: 11,
    flying: 17, bug: 18, rock: 8, ghost: 10, dragon: 14, dark: 11, steel: 10,
  },
  kalos: {
    normal: 18, fire: 8, water: 14, grass: 11, electric: 8,
    psychic: 14, ice: 5, fighting: 8, poison: 5, ground: 8,
    flying: 12, bug: 9, rock: 7, ghost: 9, dragon: 9, dark: 6, steel: 7, fairy: 15,
  },
};

const TYPE_COLORS: Record<string, string> = {
  normal: "#a0aec0", fire: "#f6ad55", water: "#63b3ed", grass: "#68d391",
  electric: "#f6e05e", psychic: "#f687b3", ice: "#76e4f7", fighting: "#c05621",
  poison: "#9f7aea", ground: "#d69e2e", flying: "#90cdf4", bug: "#a8e063",
  rock: "#b7791f", ghost: "#553c9a", dragon: "#7f9cf5", dark: "#4a5568",
  steel: "#718096", fairy: "#fbb6ce",
};

export function renderTypeDistributionChart(
  containerId: string,
  regionName: string,
): void {
  const el = document.getElementById(containerId);
  if (!el) return;

  const data = REGION_TYPE_DATA[regionName];
  if (!data) {
    el.innerHTML =
      '<p style="color:#718096;font-size:0.85rem;text-align:center;padding:1rem">Sin datos de distribución</p>';
    return;
  }

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

  window.addEventListener("resize", () => chart.resize());
}
