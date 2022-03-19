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

export const ActionTable = () => {
    return (
        <MaterialTable
            title="Table with actions"
            actions={[
                {
                    icon: tableIcons.Delete,
                    tooltip: "Delete User",
                    onClick: (event, rowData) => alert("You want to delete " + rowData.name),
                },
                {
                    icon: tableIcons.Add,
                    tooltip: "Add User",
                    isFreeAction: true,
                    onClick: (event) => alert("You want to add a new row"),
                },
            ]}
            icons={tableIcons}
            columns={columns}
            data={data}
        />
    );
};
