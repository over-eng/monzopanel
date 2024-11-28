import React from 'react';
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from 'recharts';
import styles from "./OvertimeLineChart.module.css"

export interface DataPoint {
  name: string;
  value: number;
}

const GradientLineChart = ({ data }: { data: DataPoint[] }) => {
  const renderDot = (props: any): React.ReactElement<SVGElement> => {
    const { cx, cy, index } = props;
    
    // Return with fragment for dots we don't want to render
    if (index !== 0 && index !== data.length - 1) {
      return <></>;
    }
    
    return (
      <circle
        key={`dot-${index}`}
        cx={cx}
        cy={cy}
        r={15}
        fill="#FF7F50"
        stroke="white"
        strokeWidth={2}
      />
    );
  };

  const dateFormatter = (tickItem: string | number): string => {
    const date = new Date(tickItem);
    return `${date.getDate()}/${date.getMonth() + 1}`;
  };

  return (
    <div className={styles.container}>
      <ResponsiveContainer width="100%" height="100%">
        <AreaChart
          data={data} 
          margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
        >
          <defs key="gradient">
            <linearGradient id="colorGradient" x1="0" y1="0" x2="0" y2="1">
              <stop offset="5%" stopColor="#FF7F50" stopOpacity={0.8}/>
              <stop offset="95%" stopColor="#FF7F50" stopOpacity={0}/>
            </linearGradient>
          </defs>
          <CartesianGrid key="grid" strokeDasharray="3 3" />
          <XAxis
            key="x-axis"
            dataKey="name"
            tickFormatter={dateFormatter}
            interval={24}
          />
          <YAxis key="y-axis"/>
          <Tooltip key="tooltip" />
          <Area
            key="total-events"
            type="monotone"
            dataKey="value"
            stroke="#FF7F50"
            fill="url(#colorGradient)"
            dot={renderDot}
            activeDot={{ r: 6, fill: "#FF7F50", stroke: "white", strokeWidth: 2 }}
            isAnimationActive={false}
          />
        </AreaChart>
      </ResponsiveContainer>
    </div>
  );
};

export default GradientLineChart;