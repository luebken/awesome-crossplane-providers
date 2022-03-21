import "./Footer.css";
import data from "./data";


function Footer() {
    return (
        <div className="Footer">
            Generated: {data.date}
        </div>
    );
}

export default Footer;
