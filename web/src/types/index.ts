export interface Station {
  ID: string
  Slug: string
  DisplayName: string
  Description: string
  Genre: string
  WebsiteURL: string | null
  ArtKey: string | null
  Status: 'active' | 'suspended' | 'revoked'
  CreatedAt: string
  UpdatedAt: string
}

export interface Application {
  ID: string
  SessionID: string
  ContactEmail: string
  StationName: string
  Description: string
  Genre: string
  Notes: string
  Status: 'pending' | 'approved' | 'denied'
  ReviewedBy: string | null
  ReviewedAt: string | null
  StationID: string | null
  CreatedAt: string
}

export interface AdminAction {
  ID: string
  AdminID: string
  Action: string
  TargetType: string
  TargetID: string
  Reason: string
  CreatedAt: string
}

export interface GuestInfo {
  id: string
  name: string
  muted_by_host: boolean
  muted_self: boolean
}

export interface RoomInfo {
  room_id: string
  guest_id: string
  station_id: string
  guests: GuestInfo[]
}

export interface ApprovalResult {
  station_id: string
  station_key: string
  access_token: string
}

export type Page = 'home' | 'directory' | 'apply' | 'admin' | 'admin-dashboard'
