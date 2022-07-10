// swiftlint:disable all
import Foundation
struct Translations {
	static let SOMETHING = NSLocalizedString("SOMETHING", comment: "")
	static let SOMETHING_ESCAPED = NSLocalizedString("SOMETHING_ESCAPED", comment: "")
	static func SOMETHING_WITH_ARGUMENTS(_ p1: String, _ p2: String) -> String { return NSLocalizedString("SOMETHING_WITH_ARGUMENTS", comment: "").replacingOccurrences(of: "%1", with: p1).replacingOccurrences(of: "%2", with: p2) }
}
