import MaterialTable from "material-table";
import tableIcons from "./MaterialTableIcons";

const data = [
    {
        name: "Mohammad",
        surname: "Faisal",
        birthYear: 1995,
        imageUrl: "https://avatars0.githubusercontent.com/u/7895451?s=460&v=4",
    },
    {
        name: "Nayeem Raihan ",
        surname: "Shuvo",
        birthYear: 1994,
        imageUrl: "https://avatars0.githubusercontent.com/u/7895451?s=460&v=4",
    },
];

const columns = [
    {
        title: "Avatar",
        field: "imageUrl",
        render: (rowData) => <img src={rowData.imageUrl} style={{ width: 40, borderRadius: "50%" }} />,
    },
    { title: "Name", field: "name" },
    { title: "Surname", field: "surname" },
    { title: "Birth Year", field: "birthYear", type: "numeric" },
];

export const ImageTable = () => {
    return <MaterialTable title="Image Table" icons={tableIcons} columns={columns} data={data} />;
};
