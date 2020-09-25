import { Card } from "antd";
import moment from "moment";
import React from "react";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { NSHandler } from "shared/components";
import useAPI from "shared/hooks";
import "./styles.scss";

const axisPros = { axisLine: false, tickSize: 0, tickMargin: 8 };

const formatXAxis = (dateStr) => moment(dateStr).format("Do");

const formatTooltip = (v) => [v, "customers"];

function Home() {
  const [newCustomersDaily = [], status] = useAPI("/api/v1/customers/aggregations/new");
  return (
    <div className="home">
      <Card title="New Customers" size="small">
        <NSHandler status={status}>
          {() => (
            <ResponsiveContainer height={300}>
              <AreaChart data={newCustomersDaily}>
                <defs>
                  <linearGradient id="uvcolor" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#00cc83" stopOpacity={0.95} />
                    <stop offset="95%" stopColor="#FFFFFF" stopOpacity={0.1} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#dddddd" />
                <Area type="monotone" dataKey="count" stroke="#00cc83" strokeWidth={2} fill="url(#uvcolor)" />
                <Tooltip formatter={formatTooltip} labelFormatter={formatXAxis} />
                <XAxis dataKey="date" {...axisPros} tickFormatter={formatXAxis} />
                <YAxis {...axisPros} width={32} />
              </AreaChart>
            </ResponsiveContainer>
          )}
        </NSHandler>
      </Card>
    </div>
  );
}

export default Home;
