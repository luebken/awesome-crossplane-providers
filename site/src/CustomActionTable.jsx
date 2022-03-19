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

export const CustomActionTable = () => {
    return (
        <MaterialTable
            title="Custom Action Table"
            actions={[
                {
                    icon: tableIcons.Delete,
                    tooltip: "Delete User",
                    onClick: (event, rowData) => alert("You want to delete " + rowData.name),
                },
            ]}
            components={{
                Action: (props) => (
                    <button onClick={(event) => props.action.onClick(event, props.data)}>Custom Delete Button</button>
                ),
            }}
            icons={tableIcons}
            columns={columns}
            data={data}
        />
    );
};
