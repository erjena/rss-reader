import React from 'react';
import '../public/main.css';

export default class Feed extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    const list = this.props.data.map((item, index) =>
    <div className="title">
      <a href={item.link} key={index} style={{ color: "black", textDecoration: "none" }}>
        {item.title}
      </a>
    </div>
    )
    return (
      <div>
        {list}
      </div>
    )
  }
}
