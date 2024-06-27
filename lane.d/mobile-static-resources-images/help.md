NAME
```
    mobile-static-resources-images
    - a lane action
```

SYNOPSIS
```
    -i directory -o file
    -h
```

DESCRIPTION
```
    DEPRECATED - this functionality is now built into Xcode 15 and newer. See https://developer.apple.com/documentation/xcode-release-notes/xcode-15-release-notes#Asset-Catalogs

    Searches an assets directory for .imagesets and generates swift code.

    The purpose is to enable compilation checks for image references to embedded assets.
```

OPTIONS
```
    -h
        Shows the full help.

    -i
        A directory to search for .imageset in.

    -o
        A path for the generated output file.
```

EXAMPLES

If `Assets.xcassets` contains a single .imageset named test.

The output using `-i Assets.xcassets -o -` would be:

```swift
    // swiftlint:disable all
    import UIKit
    struct Images {
        static let test = UIImage(named:"test")!
    }
```
