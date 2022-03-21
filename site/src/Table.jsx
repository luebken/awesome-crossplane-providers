import MaterialTable from "material-table";
import tableIcons from "./MaterialTableIcons";
import data from "./data";



export const Table = () => {
    return (
        <MaterialTable
            title="Awesome Crossplane Providers"
            icons={tableIcons}
            columns={data.columns}
            data={data.data}
            options={{ sorting: true, filtering: true, pageSize: 20 }}
        />
    );
};
