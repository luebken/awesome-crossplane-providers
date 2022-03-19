import "./App.css";
import { BasicTable } from "./BasicTable";
import { ActionTable } from "./ActionTable";
import { CustomActionTable } from "./CustomActionTable";
import { ImageTable } from "./ImageTable";
import { ExportTable } from "./ExportTable";
import { GroupTable } from "./GroupTable";

function App() {
    return (
        <div className="App">
            <BasicTable />
            <ActionTable />
            <CustomActionTable />
            <ImageTable />
            <ExportTable />
            <GroupTable />
        </div>
    );
}

export default App;
