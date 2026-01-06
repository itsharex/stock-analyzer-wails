import React from 'react'
import { IndustryInfo } from '../types'
import { Building2, Tag } from 'lucide-react'

interface IndustryPanelProps {
  industryInfo: IndustryInfo | undefined
}

const IndustryPanel: React.FC<IndustryPanelProps> = ({ industryInfo }) => {
  if (!industryInfo) {
    return (
      <div className="bg-white rounded-lg shadow-md p-4 col-span-1">
        <h3 className="text-lg font-semibold text-slate-800 mb-2 flex items-center">
          <Building2 className="w-5 h-5 mr-2 text-amber-500" /> 行业宏观
        </h3>
        <p className="text-sm text-slate-500">数据加载中...</p>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-4 col-span-1">
      <h3 className="text-lg font-semibold text-slate-800 mb-2 flex items-center">
        <Building2 className="w-5 h-5 mr-2 text-amber-500" /> 行业宏观
      </h3>
      <div className="space-y-3 text-sm">
        <div className="flex items-center">
          <span className="text-slate-500 w-20">所属行业:</span>
          <span className="font-medium text-slate-700">{industryInfo.industryName}</span>
        </div>
        <div className="flex items-center">
          <span className="text-slate-500 w-20">行业排名:</span>
          <span className="font-medium text-slate-700">{industryInfo.industryRank}</span>
        </div>
        <div className="flex items-start">
          <span className="text-slate-500 w-20 mt-1">概念标签:</span>
          <div className="flex flex-wrap gap-1">
            {industryInfo.conceptTags.map((tag, index) => (
              <span key={index} className="px-2 py-0.5 text-xs bg-indigo-100 text-indigo-700 rounded-full flex items-center">
                <Tag className="w-3 h-3 mr-1" /> {tag}
              </span>
            ))}
          </div>
        </div>
        <div className="flex items-start">
          <span className="text-slate-500 w-20 mt-1">政策影响:</span>
          <span className="font-medium text-slate-700">{industryInfo.policyImpact}</span>
        </div>
      </div>
    </div>
  )
}

export default IndustryPanel
