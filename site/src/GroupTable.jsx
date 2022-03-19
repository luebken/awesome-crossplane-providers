import MaterialTable from "material-table";
import tableIcons from "./MaterialTableIcons";

const data = [
    { birthYear: 1995, name: "Mohammad", surname: "Faisal" },
    { birthYear: 1996, name: "Mohammad", surname: "Kashem" },
    { birthYear: 1997, name: '"Nayeem Raihan', surname: "Shuvo" },
];

const columns = [
    { title: "Name", field: "name" },
    { title: "Surname", field: "surname" },
    { title: "Birth Year", field: "birthYear", type: "numeric" },
];

export const GroupTable = () => {
    return (
        <MaterialTable
            title="Group Table"
            icons={tableIcons}
            columns={columns}
            data={data}
            options={{
                grouping: true,
            }}
        />
    );
};
