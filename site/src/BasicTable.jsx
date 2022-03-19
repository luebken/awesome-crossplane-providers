import MaterialTable from "material-table";
import tableIcons from "./MaterialTableIcons";

const data = [
    { name: "Mohammad", surname: "Faisal", birthYear: 1995 },
    { name: "Nayeem Raihan ", surname: "Shuvo", birthYear: 1994 },
];

const columns = [
    { title: "Name", field: "name" },
    { title: "Surname", field: "surname" },
    { title: "Birth Year", field: "birthYear", type: "numeric" },
];

export const BasicTable = () => {
    return (
        <MaterialTable
            title="Awesome Crossplane Providers"
            icons={tableIcons}
            columns={columns}
            data={data}
            options={{ sorting: true }}
        />
    );
};
