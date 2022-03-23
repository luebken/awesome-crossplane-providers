import MaterialTable from "material-table";
import tableIcons from "./MaterialTableIcons";
import data from "./data";
import "./App.css";
//import StarRate from "@material-ui/icons/StarRate";



export const Table = () => {
    return (
        <MaterialTable
            title="Awesome Crossplane Providers"
            icons={tableIcons}
            columns={data.columns}
            data={data.data}
            detailPanel={rowData => {
                return (
                    <div className="detailPane">
                        <table className="detailTable">
                            <tbody>
                                <tr>
                                    <td className="detailTitle">{rowData.name} </td>
                                    <td></td>
                                    <td></td>
                                </tr>
                                <tr>
                                    <td> {rowData.description} </td>
                                    <td>Repository:</td>
                                    <td>CRDs:</td>
                                </tr>
                                <tr>
                                    <td></td>
                                    <td>&nbsp;&nbsp;Updated: {rowData.updated}</td>
                                    <td>&nbsp;&nbsp;Alpha: {rowData.crdsAlpha}</td>
                                </tr>
                                <tr>
                                    <td> </td>
                                    <td>&nbsp;&nbsp;Last Release: {rowData.lastReleaseDate} {rowData.lastReleaseTag}</td>
                                    <td>&nbsp;&nbsp;Beta: {rowData.crdsBeta}</td>
                                </tr>
                                <tr className="stars">
                                    <td> Stars: {rowData.stargazers}&nbsp;&nbsp;&nbsp;&nbsp;Open Issues: {rowData.openIssues}</td>
                                    <td>&nbsp;&nbsp;Created: {rowData.created}</td>
                                    <td>&nbsp;&nbsp;V1: {rowData.crdsV1}</td>
                                </tr>
                            </tbody>
                        </table>
                    </div >
                )
            }}
            options={{ sorting: true, filtering: true, pageSize: 20, pageSizeOptions: [10, 20, 50, 100] }}
        />
    );
};
