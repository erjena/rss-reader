import React from 'react';
import '../public/main.css';

export default class Sources extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    return (
      <div className="sources">
        <span style={{ paddingRight: "100px" }}>All</span>
      </div>
    )
  }
}
