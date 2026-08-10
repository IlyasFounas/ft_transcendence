import { useState } from 'react'

import OperationHistoryRow from './components/rows/OperationHistoryRow.tsx'

function App() {

  return (
	< div className="w-full min-h-screen bg-white">
      <OperationHistoryRow
        operationId="OP-004540"
        dateTime="12 JUN 2026 21:35"
        points={2025}
      />
      <OperationHistoryRow
        operationId="OP-004541"
        dateTime="12 JUN 2026 21:38"
        points={8562}
      />
	</div>
  )
}

export default App
