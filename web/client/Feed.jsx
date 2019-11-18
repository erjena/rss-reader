import React from 'react';
import '../public/main.css';

export default class Feed extends React.Component {
  constructor(props) {
    super(props)
  }

  render() {
    const list = this.props.data.map((item, index) =>
      <a href={item.link} key={index}>
        <div className="feedElement" >
          <div>
            <span className="title">
              {(item.title).replace(/<\/?[^>]+(>|$)/g, "")}
            </span>
            <span className="pubDate">
              {item.pubDate.slice(0, 17)}
            </span>
            <div style={{ clear: "both" }}></div>
          </div>
          <span className="description">{(item.description).replace(/<\/?[^>]+(>|$)/g, "")}</span>
        </div>
      </a>
    )
    return (
      <div className="feed">
        {list}
      </div>
    )
  }
}
