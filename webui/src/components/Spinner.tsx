import React from 'react';
import {HashLoader} from "react-spinners";

const Spinner : React.FunctionComponent<{text: string}> = ({text}) => {
    return (
        <div>
            <div className={""}>
                <HashLoader sizeUnit={"px"} size={30} color={'#46765E'} loading={true} css={"margin-left: auto; margin-right: auto;"}/>
            </div>
            {text.length > 0 ? <h4 className={"mt-2 font-italic"}>Loading files</h4> : <div/>}
        </div>);
};

export default Spinner;