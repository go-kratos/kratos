def _needs_install(name, dep, hkeys=["sha256", "sha1", "tag"], verbose=0, strict=False):

    # Does it already exist?
    existing_rule = native.existing_rule(name)
    if not existing_rule:
        return True

    # If it has already been defined and our dependency lists a
    # hash, do these match? If a hash mismatch is encountered, has
    # the user specifically granted permission to continue?
    for hkey in hkeys:
        expected = dep.get(hkey)
        actual = existing_rule.get(hkey)
        if expected:
            # By guarding the check of expected vs actual with a
            # pre-check for actual=None (or empty string), we're
            # basically saying "If the user did not bother to set a
            # sha256 or sha1 hash for the rule, they probably don't
            # care about overriding a dependency either, so don't
            # complain about it."  In particular, rules_go does not a
            # set a sha256 for their com_google_protobuf http_archive
            # rule, so this gives us a convenient loophole to prevent
            # collisions on this dependency.  The "strict" flag can be
            # used as as master switch to disable blowing up the
            # loading phase due to dependency collisions.
            if actual and expected != actual and strict:
                msg = """
An existing {0} rule '{1}' was already loaded with a {2} value of '{3}'.  Refusing to overwrite this with the requested value ('{4}').
Either remove the pre-existing rule from your WORKSPACE or exclude it from loading by rules_protobuf (strict={5}.
""".format(existing_rule["kind"], name, hkey, actual, expected, strict)

                fail(msg)
            else:
                if verbose > 1: print("Skip reload %s: %s = %s" % (name, hkey, actual))
                return False

    # No kheys for this rule - in this case no reload; first one loaded wins.
    if verbose > 1: print("Skipping reload of existing target %s" % name)
    return False


def _install(deps, verbose, strict):
    """Install a list if dependencies for matching native rules.
    Return:
      list of deps that have no matching native rule.
    """
    todo = []

    for d in deps:
        name = d.get("name")
        rule = d.pop("rule", None)
        if not rule:
            fail("Missing attribute 'rule': %s" % name)
        if hasattr(native, rule):
            rule = getattr(native, rule)
            if verbose: print("Loading %s)" % name)
            rule(**d)
        else:
            d["rule"] = rule
            todo.append(d)

    return todo


def require(keys,
            deps = {},
            overrides = {},
            excludes = [],
            verbose = 0,
            strict = False):

    #
    # Make a list of non-excluded required deps with merged data.
    #
    required = []

    for key in keys:
        dep = deps.get(key)
        if not dep:
            fail("Unknown workspace dependency: %s" % key)
        d = dict(**dep) # copy the 'frozen' object.
        if not key in excludes:
            over = overrides.get(key)
            data = d + over if over else d
            if _needs_install(key, data, verbose=verbose, strict=strict):
                data["name"] = key
                required.append(data)

    return _install(required, verbose, strict)