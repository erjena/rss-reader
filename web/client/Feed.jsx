import React from 'react';
import '../public/main.css';

export default class Feed extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    const list = this.props.data.map((item, index) =>
    <div key={index} className="feedElement">
      <span><a href={item.link}>
        {item.title}
      </a></span>
      <span className="pubDate">{item.pubDate.slice(0, 17)}</span>
      <span className="description">{(item.description).replace(/<\/?[^>]+(>|$)/g, "")}</span>
    </div>
    )
    return (
      <div className="feedContainer">
        {list}
      </div>
    )
  }
}
