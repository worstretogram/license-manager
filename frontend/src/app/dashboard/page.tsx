'use client'
import { MdContentCopy } from "react-icons/md";
import { BsFillTrashFill } from "react-icons/bs";
import Button from "@/components/Button";
import { useEffect, useState } from "react";
import { observer } from "mobx-react-lite";
import store from "@/store/store";
import AddModal from "@/components/AddModal";
import { deleteLicense, downloadLicense, getLicenses } from "@/api/licence";



const DashboardPage = observer(() => {
  const [isDelete, setIsDelete] = useState(false);
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

 
  useEffect(() => {
    getLicenses()
  }, [])

  const toggleDelete = () => setIsDelete((old) => !old);

  const handleCopy = (key: string) => {
    navigator.clipboard.writeText(key);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2000);
  };

  return (
    <>
      {store.isAdding && <AddModal />}
      <div className="flex flex-col gap-4">
        <div className="flex justify-end gap-4">
          <Button key="delete" onClick={toggleDelete} title="Delete license" />
          <Button key="create" onClick={() => store.changeIsAdding(true)} title="Creating a license" />
        </div>
        <table className="border-border text-white font-bold w-full">
          <thead className="text-center bg-secondary">
            <tr>
              {["Company", "License ID", "Status", "Expires", "Users", "Signature", "Download"]
                .concat(isDelete ? [""] : [])
                .map((h, idx) => (
                  <th key={h + idx} className="px-4 py-3 border-b border border-border">
                    {h}
                  </th>
                ))}
            </tr>
          </thead>
          <tbody>
            {store.licenses?.map((row, i) => (
              <tr key={row.license_id || i} className="bg-secondary border-b border-border">
                <td className="border px-4 py-3">{row.owner}</td>
                <td className="border px-4 py-3">{row.license_id}</td>
                <td className="border px-4 py-3">
                  <span className="flex items-center gap-2">
                    <span className="w-3 h-3 rounded-full bg-green-500" />
                    Active
                  </span>
                </td>
                <td className="border px-4 py-3">{row.expires_at?.split("T")[0]}</td>
                <td className="border px-4 py-3">{row.max_users}</td>
                <td className="border px-4 py-3 relative">
                  <div
                    className="flex items-center gap-2 cursor-pointer"
                    title="Click to copy"
                    onClick={() => handleCopy(row.signature)}
                  >
                    <span className="truncate max-w-[200px]">{row.signature}</span>
                    <MdContentCopy />
                  </div>
                 
                  {copiedKey === row.signature && (
                    <span className="absolute -top-6 left-1/2 -translate-x-1/2 px-2 py-0.5 text-xs text-white bg-green-600 rounded shadow transition-opacity duration-300">
                      Copied!
                    </span>
                  )}
                </td>
                <td onClick={() => downloadLicense(row.license_id)} className="border px-4 py-3">Download</td>
                {isDelete && (
                  <td className="border px-4 py-3">
                    <div onClick={() => deleteLicense(row.license_id)} className="flex justify-center items-center cursor-pointer" title="Delete">
                      <BsFillTrashFill color="#FF3B30" />
                    </div>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
});

export default DashboardPage;
