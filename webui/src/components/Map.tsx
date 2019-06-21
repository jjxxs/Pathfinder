import React from 'react';
import {connect} from "react-redux";
import {AppState} from "../redux/AppState";
import Spinner from "./Spinner";

interface MapProps {
    image: string,
    points: number[],
}

class ImageLoader {
    private image = new Image();
    private initialized = false;

    setImage(image: string) {
        if (!this.initialized) {
            this.image.src = "data:image/gif;base64," + image;
        }
        this.initialized = true;
    }

    isInitialized() : boolean {
        return this.initialized;
    }

    getImage() : HTMLImageElement {
        return this.image;
    }
}

const imageLoader = new ImageLoader();

const Canvas : React.FunctionComponent<MapProps> = ({image, points}) => {
    const img = new Image();
    img.src = "data:image/gif;base64," + image;

    const actualPoints = [];
    for (let i = 0; i < points.length; i += 2) {
        actualPoints.push({X: points[i], Y: points[i+1]});
    }

    const scaling = 620 / img.height;

    const lines = [];
    for (let i = 0; i < actualPoints.length; i++) {
        if (i === actualPoints.length-1) {
            let x1 = actualPoints[i].X * scaling;
            let x2 = actualPoints[0].X * scaling;
            let y1 = actualPoints[i].Y * scaling;
            let y2 = actualPoints[0].Y * scaling;
            lines.push(<line x1={x1} x2={x2} y1={y1} y2={y2} stroke={"red"} strokeWidth={"1"}/>)
        } else {
            let x1 = actualPoints[i].X * scaling;
            let x2 = actualPoints[i+1].X * scaling;
            let y1 = actualPoints[i].Y * scaling;
            let y2 = actualPoints[i+1].Y * scaling;
            lines.push(<line x1={x1} x2={x2} y1={y1} y2={y2} stroke={"red"} strokeWidth={"1"}/>)
        }
    }

    return (
        <div className={"m-auto j-image-container"}>
            <div className={"col j-image"}>
                <svg width={img.width * scaling} height={620} xmlns="http://www.w3.org/2000/svg" version="1.1">
                    {lines}
                </svg>
                <img src={"data:image/gif;base64," + image} alt={"Not available."}/>
              </div>
           </div>);
};

const CanvasContainer : React.FunctionComponent<MapProps> = ({image, points}) => {
    const actualImage = image.length === 0 ?
        <div className={"m-auto"}><Spinner text={""}/></div> :
        <Canvas image={image} points={points}/>;

    if (image.length > 0 && !imageLoader.isInitialized()) {
        imageLoader.setImage(image);
    }

    return (
        <div className={"col d-flex map"}>
            {actualImage}
        </div>);
};

const mapStateToProps = (state: AppState) => {
    return {image: state.image, points: state.points}
};

export default connect(mapStateToProps)(CanvasContainer);