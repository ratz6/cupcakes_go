import React, { Component } from "react";
// useEffect, useState, PureComponent, 
import {
  LineChart,
  Line,
  Tooltip,
  CartesianGrid,
  XAxis,
  YAxis
} from "recharts";
import axios from "axios";

import "./styles.css";

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      res: []
    };
  }

  interval = null;
  componentDidMount() {
    this.interval = setInterval(this.getData, 30000);
    this.getData();
  }

  componentWillUnmount() {
    clearInterval(this.interval);
  }

  getData = () => {
    axios
      .get("API LINK", {
        responseType: "json"
      })
      .then((response) => {
        this.setState({ res: response.data });
      });
  };

  render() {
    const { res } = this.state;
    return (
      <LineChart
        width={1300}
        height={300}
        data={res}
        margin={{
          top: 10,
          right: 10,
          left: 10,
          bottom: 5
        }}
      >
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="month" />
        <YAxis />
        <Tooltip />
        {/* <Legend /> */}
        <Line type="monotone" dataKey="cupcake" stroke="#8884d8" />
      </LineChart>
    );
  }
}
export default App;