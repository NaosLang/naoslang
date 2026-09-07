module.exports = grammar({
    name: 'naoslang',

    extras: $ => [
        /\s/,
        /\/\/.*/,
    ],

    rules: {
        source_file: $ => seq(
            repeat($._imports),
        ),

        // +---------+
        // | Imports |
        // +---------+

        _imports: $ => choice(
            $.global_import,
            $.alias_import,
        ),

        _generic_import: $ => seq(
            '@import',
            '(',
            field('path', $.string_literal),
            ')',
        ),

        global_import: $ => seq(
            'using',
            $._generic_import,
            ';',
        ),

        alias_import: $ => seq(
            $.identifier,
            '=',
            $._generic_import,
            ';',
        ),

        // +------+
        // | Misc |
        // +------+

        string_literal: $ => /"([^"\\]|\\.)*"/,
        identifier: $ => /[a-zA-Z_][a-zA-Z0-9_]*/,
    }
})





